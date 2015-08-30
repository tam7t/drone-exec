package yaml

import "gopkg.in/yaml.v2"

// Parse parses a Yaml configuraiton file.
func Parse(raw []byte) (*Config, error) {
	c := Config{}
	e := yaml.Unmarshal(raw, &c)
	return &c, e
}

// ParseString parses a Yaml configuration file
// in string format.
func ParseString(raw string) (*Config, error) {
	return Parse([]byte(raw))
}
