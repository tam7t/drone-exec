package rules

import (
	"strings"

	"github.com/drone/drone-exec/yaml"
)

type StepFunc func(*yaml.Step) error

// forEachStep is a helper function that executes
// a particular rule across all steps in the
// build process.
func forEachStep(c *yaml.Config, fn StepFunc) error {
	var steps []*yaml.Step
	steps = append(steps, c.Cache)
	steps = append(steps, c.Clone)
	steps = append(steps, c.Build)

	for _, step := range c.Publish {
		steps = append(steps, step)
	}
	for _, step := range c.Deploy {
		steps = append(steps, step)
	}
	for _, step := range c.Notify {
		steps = append(steps, step)
	}
	for _, step := range c.Compose {
		steps = append(steps, step)
	}

	for _, step := range steps {
		if step == nil {
			continue
		}
		err := fn(step)
		if err != nil {
			return err
		}
	}

	return nil
}

// pluginName is a helper function that resolves the
// plugin image name. When using official drone plugins
// it is possible to use an alias name. This converts to
// the fully qualified name.
func pluginName(name string) string {
	if strings.Contains(name, "/") {
		return name
	}
	name = strings.Replace(name, "_", "-", -1)
	name = "plugins/drone-" + name
	return name
}

// pluginNameDefault is a helper function that resolves
// the plugin image name. If the image name is blank the
// default name is used instead.
func pluginNameDefault(name, defaultName string) string {
	if len(name) == 0 {
		name = defaultName
	}
	return pluginName(name)
}
