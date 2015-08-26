package yaml

import (
	"gopkg.in/yaml.v2"
)

func Parse(raw string) (*Config, error) {
	conf := &Config{}
	err := yaml.Unmarshal([]byte(raw), conf)
	return conf, err
}
