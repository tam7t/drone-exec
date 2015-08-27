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

// prepareClone prepares the clone object. It applies
// default settings if none exist.
func prepareClone(c *yaml.Config) error {
	if c.Clone == nil {
		c.Clone = &yaml.Step{}
	}
	c.Clone.Image = pluginNameDefault(
		c.Clone.Image,
		DefaultCloner,
	)
	return nil
}

// prepareCache prepares the cache object. It applies
// default settings if none exist.
func prepareCache(c *yaml.Config, name string) error {
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

// prepareImages prepares all images names.
func prepareImages(c *yaml.Config) error {
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

// verifyCache verifies the cache section of the yaml
// is setup correctly.
func verifyCache(c *yaml.Config) error {
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

// verifyBuild verifies the build section of the yaml
// is present and has a valid image name.
func verifyBuild(c *yaml.Config) error {
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

// verifyPlugins verifies the plugins are part of the
// plugin white-list.
func verifyPlugins(c *yaml.Config, match string) error {
	// always use the default plugin filter if no
	// matching string is provided. Safety first!
	if len(match) == 0 {
		match = DefaultMatch
	}

	return forEachStep(c, func(s *yaml.Step) error {
		// the build step is not a plugin, and therefore
		// is not subject to the plugin whitelist.
		if s == c.Build {
			return nil
		}

		// verify the user specified the plugin image.
		if len(s.Image) == 0 {
			return fmt.Errorf("Yaml must define plugin images")
		}

		// uses filepath globbing for plugin matching
		ok, _ := filepath.Match(match, s.Image)
		if ok {
			return nil
		}
		return fmt.Errorf("Yaml image %s is forbidden", s.Image)
	})
}

// verifyPlugins verifies all build steps have associated
// docker images defined.
func verifyImages(c *yaml.Config) error {
	return forEachStep(c, func(s *yaml.Step) error {
		if len(s.Image) != 0 {
			return nil
		}
		return fmt.Errorf("Yaml must define an image for every step")
	})
}

// cleanVolumes is a rule that ensures every
// step is executed without volumes.
func cleanVolumes(c *yaml.Config) error {
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

// cleanNetwork is a transformer that ensures every
// step is executed with default bridge networking.
func cleanNetwork(c *yaml.Config) error {
	return forEachStep(c, func(s *yaml.Step) error {
		s.NetworkMode = ""
		return nil
	})
}

// cleanPrivileged is a transformer that ensures every
// step is executed in non-privileged mode.
func cleanPrivileged(c *yaml.Config) error {
	forEachStep(c, func(s *yaml.Step) error {
		s.Privileged = false
		return nil
	})

	// the only white-listed plugin that can
	// run in privileged mode is the `drone-docker`
	// plugin.
	for _, step := range c.Publish {
		if step.Image == "plugins/drone-docker" {
			step.Privileged = true
			step.Volumes = nil
			step.NetworkMode = ""
			step.Entrypoint = []string{}
		}
	}
	return nil
}
