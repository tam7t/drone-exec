package yaml

import "gopkg.in/yaml.v2"

// Parse parses a Yaml configuraiton file.
func Parse(in []byte) (*Config, error) {
	c := Config{}
	e := yaml.Unmarshal(in, &c)
	return &c, e
}

// ParseString parses a Yaml configuration file
// in string format.
func ParseString(in string) (*Config, error) {
	return Parse([]byte(in))
}
