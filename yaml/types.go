package yaml

// Config is a typed representation of the
// Yaml configuration file.
type Config struct {
	Cache Plugin
	Clone Plugin
	Build Build

	Compose Containerslice
	Publish Pluginslice
	Deploy  Pluginslice
	Notify  Pluginslice
}

// Container is a typed representation of a
// docker step in the Yaml configuration file.
type Container struct {
	Image            string
	Pull             bool
	Privileged       bool
	DisableAptConfig bool `yaml:"disable_apt_config"`
	Environment      MapEqualSlice
	Entrypoint       Command
	Command          Command
	ExtraHosts       []string `yaml:"extra_hosts"`
	Volumes          []string
	Net              string
	AuthConfig       AuthConfig `yaml:"auth_config"`
}

// Build is a typed representation of the build
// step in the Yaml configuration file.
type Build struct {
	Container `yaml:",inline"`

	Commands []string
}

// Auth for Docker Image Registry
type AuthConfig struct {
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Email         string `yaml:"email"`
	RegistryToken string `yaml:"registry_token"`
}

// Plugin is a typed representation of a
// docker plugin step in the Yaml configuration
// file.
type Plugin struct {
	Container `yaml:",inline"`

	Vargs  Vargs  `yaml:",inline"`
	Filter Filter `yaml:"when"`
}

// Vargs holds unstructured arguments, specific
// to the plugin, that are used at runtime when
// executing the plugin.
type Vargs map[string]interface{}

// Filter is a typed representation of filters
// used at runtime to decide if a particular
// plugin should be executed or skipped.
type Filter struct {
	Repo    string
	Branch  Stringorslice
	Event   Stringorslice
	Success string
	Failure string
	Change  string
	Matrix  map[string]string
}
