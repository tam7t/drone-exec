package yaml

import (
	"fmt"
	"strings"

	"github.com/flynn/go-shlex"
	"gopkg.in/yaml.v2"
)

type Command struct {
	parts []string
}

func (s *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stringType string
	err := unmarshal(&stringType)
	if err == nil {
		s.parts, err = shlex.Split(stringType)
		return err
	}

	var sliceType []string
	err = unmarshal(&sliceType)
	if err == nil {
		s.parts = sliceType
		return nil
	}

	return err
}

func (s *Command) Slice() []string {
	return s.parts
}

type MapEqualSlice struct {
	parts []string
}

func (s *MapEqualSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&s.parts)
	if err == nil {
		return nil
	}

	var mapType map[string]string

	err = unmarshal(&mapType)
	if err != nil {
		return err
	}

	for k, v := range mapType {
		s.parts = append(s.parts, strings.Join([]string{k, v}, "="))
	}

	return nil
}

func (s *MapEqualSlice) Slice() []string {
	return s.parts
}

// Pluginslice is a slice of Plugins with a custom Yaml
// unarmshal function to preserve ordering.
type Pluginslice struct {
	parts []*Plugin
}

func (s *Pluginslice) UnmarshalYAML(unmarshal func(interface{}) error) error {

	// unmarshal the yaml into the generic
	// mapSlice type to preserve ordering.
	obj := yaml.MapSlice{}
	err := unmarshal(&obj)
	if err != nil {
		return err
	}

	// unarmshals each item in the mapSlice,
	// unmarshal and append to the slice.
	err = unmarshalYaml(obj, func(key string, val []byte) error {
		plugin := &Plugin{}
		err := yaml.Unmarshal(val, plugin)
		if err != nil {
			return err
		}
		if len(plugin.Image) == 0 {
			plugin.Image = key
		}
		s.parts = append(s.parts, plugin)
		return nil
	})
	return err
}

func (s *Pluginslice) Slice() []*Plugin {
	return s.parts
}

// ContainerSlice is a slice of Containers with a custom
// Yaml unarmshal function to preserve ordering.
type Containerslice struct {
	parts []*Container
}

func (s *Containerslice) UnmarshalYAML(unmarshal func(interface{}) error) error {

	// unmarshal the yaml into the generic
	// mapSlice type to preserve ordering.
	obj := yaml.MapSlice{}
	err := unmarshal(&obj)
	if err != nil {
		return err
	}

	// unarmshals each item in the mapSlice,
	// unmarshal and append to the slice.
	return unmarshalYaml(obj, func(key string, val []byte) error {
		ctr := &Container{}
		err := yaml.Unmarshal(val, ctr)
		if err != nil {
			return err
		}
		if len(ctr.Image) == 0 {
			ctr.Image = key
		}
		s.parts = append(s.parts, ctr)
		return nil
	})
}

func (s *Containerslice) Slice() []*Container {
	return s.parts
}

// emitter defines the callback function used for
// generic yaml parsing. It emits back a raw byte
// slice for custom unmarshalling into a structure.
type unmarshal func(string, []byte) error

// unmarshalYaml is a helper function that removes
// some of the boilerplate from unmarshalling
// complex map slices.
func unmarshalYaml(v yaml.MapSlice, emit unmarshal) error {
	for _, vv := range v {
		// re-marshal the interface{} back to
		// a raw yaml value
		val, err := yaml.Marshal(&vv.Value)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%v", vv.Key)

		// unmarshal the raw value using the
		// callback function.
		if err := emit(key, val); err != nil {
			return err
		}
	}
	return nil
}
