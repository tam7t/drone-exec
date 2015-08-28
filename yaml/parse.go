package yaml

import "gopkg.in/yaml.v2"

// RuleFunc to extend yaml parsing to validate and
// transform the results.
type RuleFunc func(*Config) error

// Config represents a repository build configuration.
type Config struct {
	Cache *Step
	Clone *Step
	Build *Step

	Compose map[string]*Step
	Publish map[string]*Step
	Deploy  map[string]*Step
	Notify  map[string]*Step
}

// Parse parses a yaml file.
func Parse(raw string) (*Config, error) {
	conf := &Config{}
	err := yaml.Unmarshal([]byte(raw), conf)
	return conf, err
}

// Parse parses a yaml file and then applies the
// transform and validation rules.
func ParseRules(raw string, rules []RuleFunc) (*Config, error) {
	conf, err := Parse(raw)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if err := rule(conf); err != nil {
			return nil, err
		}
	}
	return conf, nil
}

// yaml
//
// 1. unit test rules
// 2. rule builder ... rules.RuleSet ... script.RuleSet ... send to yaml when parsing
//
// 4. calculate workspace
//		3.a rule to calculate clone path / workspace
//		3.b how to propogate to plugins?

// runner
//
// 0. add code to docker seciton, unit test
// 1. add code to apply the rules
// 2. add code to setup the pod
// 3. add code to run the build
// 4. add code to cancel / timeout the build

// book keeping
//
// 1. implement vendored packages

// deployment
//
// 1. setup automated docker deploy
// 2. add build instructions to README
// 3. add run instructions to README
