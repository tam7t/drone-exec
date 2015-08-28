package rules

import (
	"fmt"
	"path/filepath"

	"github.com/drone/drone-exec/yaml"
)

// Default clone plugin.
const DefaultCloner = "plugins/drone-git"

// Default cache plugin.
const DefaultCacher = "plugins/drone-cache"

// Default plugin whitelist match string.
const DefaultMatch = "plugins/*"

// RuleFunc to extend yaml parsing to validate and
// transform the results.
type RuleFunc func(*yaml.Config) error

// PrepareClone prepares the clone object. It applies
// default settings if none exist.
func PrepareClone(c *yaml.Config) error {
	if c.Clone == nil {
		c.Clone = &yaml.Step{}
	}
	c.Clone.Image = pluginNameDefault(
		c.Clone.Image,
		DefaultCloner,
	)
	return nil
}

// PrepareCache prepares the cache object. It applies
// default settings if none exist.
func PrepareCache(c *yaml.Config, name string) error {
	if c.Cache == nil {
		return nil
	}
	c.Cache.Image = pluginNameDefault(
		c.Cache.Image,
		DefaultCacher,
	)
	c.Cache.Volumes = []string{
		filepath.Join("/tmp/drone/", name),
	}
	return nil
}

// PrepareCacheRule is an adapter function that allows
// PrepareCache to be executed as a RuleFunc.
func PrepareCacheRule(name string) RuleFunc {
	return func(c *yaml.Config) error {
		return PrepareCache(c, name)
	}
}

// PrepareImages prepares all images names.
func PrepareImages(c *yaml.Config) error {
	return forEachStep(c, func(s *yaml.Step) error {
		if len(s.Image) == 0 {
			return nil
		}
		if s == c.Build {
			return nil
		}
		s.Image = pluginName(s.Image)
		return nil
	})
}

// PrepareEnv prepares the default environment variables
// fo for the build image.
func PrepareEnv(c *yaml.Config, envs []string) error {
	if c.Build == nil {
		return nil
	}
	c.Build.Environment = append(c.Build.Environment, envs...)
	return nil
}

// PrepareEnvRule is an adapter function that allows
// PrepareEnv to be executed as a RuleFunc.
func PrepareEnvRule(envs []string) RuleFunc {
	return func(c *yaml.Config) error {
		return PrepareEnv(c, envs)
	}
}

// LintCache verifies the cache section of the yaml
// is setup correctly.
func LintCache(c *yaml.Config) error {
	if c.Cache == nil {
		return nil
	}
	if c.Cache.Config == nil {
		return fmt.Errorf("Yaml must define cache mountpoints")
	}
	mountv, ok := c.Cache.Config["mount"]
	if !ok {
		return fmt.Errorf("Yaml must define cache mountpoints")
	}
	_, ok = mountv.([]string)
	if !ok {
		return fmt.Errorf("Yaml has a malformed cache section")
	}
	return nil
}

// LintBuild verifies the build section of the yaml
// is present and has a valid image name.
func LintBuild(c *yaml.Config) error {
	if c.Build == nil {
		return fmt.Errorf("Yaml must define a build section")
	}
	if len(c.Build.Image) == 0 {
		return fmt.Errorf("Yaml must define a build immage")
	}
	if c.Build.Config == nil || c.Build.Config["commands"] == nil {
		return fmt.Errorf("Yaml must define build commands")
	}
	return nil
}

// LintPlugins verifies the plugins are part of the
// plugin white-list.
func LintPlugins(c *yaml.Config, patterns []string) error {
	// always use the default plugin filter if no
	// matching string is provided. Safety first!
	if len(patterns) == 0 {
		patterns = []string{DefaultMatch}
	}

	return forEachStep(c, func(s *yaml.Step) error {
		// the build step is not a plugin, and therefore
		// is not subject to the plugin whitelist.
		if s == c.Build {
			return nil
		}

		match := false
		for _, pattern := range patterns {
			if pattern == s.Image {
				match = true
				break
			}
			ok, err := filepath.Match(pattern, s.Image)
			if ok && err == nil {
				match = true
				break
			}
		}
		if !match {
			return fmt.Errorf("Yaml cannot use un-trusted image %s", s.Image)
		}
		return nil
	})
}

// LintPluginsRule is an adapter function that allows
// LintPlugins to be executed as a RuleFunc.
func LintPluginsRule(patterns []string) RuleFunc {
	return func(c *yaml.Config) error {
		return LintPlugins(c, patterns)
	}
}

// LintImages verifies all build steps have associated
// docker images defined.
func LintImages(c *yaml.Config) error {
	return forEachStep(c, func(s *yaml.Step) error {
		if len(s.Image) != 0 {
			return nil
		}
		return fmt.Errorf("Yaml must define an image for every step")
	})
}

// CleanVolumes is a rule that ensures every
// step is executed without volumes.
func CleanVolumes(c *yaml.Config) error {
	return forEachStep(c, func(s *yaml.Step) error {
		if s == c.Cache {
			// the cache plugins volumes were already
			// set and overriden to the appropriate values
			return nil
		}
		s.Volumes = nil
		return nil
	})
}

// CleanNetwork is a transformer that ensures every
// step is executed with default bridge networking.
func CleanNetwork(c *yaml.Config) error {
	return forEachStep(c, func(s *yaml.Step) error {
		s.NetworkMode = ""
		return nil
	})
}

// CleanPrivileged is a transformer that ensures every
// step is executed in non-privileged mode.
func CleanPrivileged(c *yaml.Config) error {
	return forEachStep(c, func(s *yaml.Step) error {
		s.Privileged = false
		return nil
	})
}

// EnableDocker is a transformer that ensures the Docker
// plugin is executed in privileged mode.
func EnableDocker(c *yaml.Config) error {
	for _, step := range c.Publish {
		if step.Image == "plugins/drone-docker" {
			step.Privileged = true

			// since we are running in privileged mode
			// we make sure there are no volumes, host network
			// or anything else that could be exploited
			step.Volumes = nil
			step.NetworkMode = ""
			step.Entrypoint = []string{}
		}
	}
	return nil
}
