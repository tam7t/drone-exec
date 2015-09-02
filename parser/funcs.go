package parser

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

var (
	ErrImageMissing   = errors.New("Yaml must specify an image for every step")
	ErrImageWhitelist = errors.New("Yaml must specify am image from the white-list")
)

const (
	DefaultCloner = "plugins/drone-git"   // Default clone plugin.
	DefaultCacher = "plugins/drone-cache" // Default cache plugin.
	DefaultMatch  = "plugins/*"           // Default plugin whitelist.
)

// RuleFunc defines a function used to validate or modify the yaml during
// the parsing process.
type RuleFunc func(Node) error

// ImageName expands to a fully qualified image name. If no image name is found,
// a default is used when possible, else ErrImageMissing is returned.
func ImageName(n Node) error {
	d, ok := n.(*DockerNode)
	if !ok {
		return nil
	}

	switch d.NodeType {
	case NodeBuild, NodeCompose:
		break
	case NodeClone:
		d.Image = expandImageDefault(d.Image, DefaultCloner)
	case NodeCache:
		d.Image = expandImageDefault(d.Image, DefaultCacher)
	default:
		d.Image = expandImage(d.Image)
	}

	if len(d.Image) == 0 {
		return ErrImageMissing
	}
	d.Image = expandImageTag(d.Image)
	return nil
}

// ImageMatch checks the image name against a whitelist.
func ImageMatch(n Node, patterns []string) error {
	d, ok := n.(*DockerNode)
	if !ok {
		return nil
	}
	switch d.NodeType {
	case NodeBuild, NodeCompose:
		return nil
	}
	if len(patterns) == 0 {
		patterns = []string{DefaultMatch}
	}
	match := false
	for _, pattern := range patterns {
		if pattern == d.Image {
			match = true
			break
		}
		ok, err := filepath.Match(pattern, d.Image)
		if ok && err == nil {
			match = true
			break
		}
	}
	if !match {
		return fmt.Errorf("Plugin %s is not in the whitelist", d.Image)
	}
	return nil
}

func ImageMatchFunc(patterns []string) RuleFunc {
	return func(n Node) error {
		return ImageMatch(n, patterns)
	}
}

// Sanitize sanitizes a Docker Node by removing any potentially
// harmful configuration options.
func Sanitize(n Node) error {
	d, ok := n.(*DockerNode)
	if !ok {
		return nil
	}
	d.Privileged = false
	d.Volumes = nil
	d.Net = ""
	d.Entrypoint = []string{}
	return nil
}

func SanitizeFunc(trusted bool) RuleFunc {
	return func(n Node) error {
		if !trusted {
			return Sanitize(n)
		}
		return nil
	}
}

// Escalate escalates a Docker Node to run in privileged mode if
// the plugin is whitelisted.
func Escalate(n Node) error {
	d, ok := n.(*DockerNode)
	if !ok {
		return nil
	}
	image := strings.Split(d.Image, ":")
	if d.NodeType == NodePublish &&
		image[0] == "plugins/drone-docker" {

		d.Privileged = true
		d.Volumes = nil
		d.Net = ""
		d.Entrypoint = []string{}
	}
	return nil
}

// Cache transforms the Docker Node to mount a volume to the host
// machines local cache.
func Cache(n Node, dir string) error {
	d, ok := n.(*DockerNode)
	if !ok {
		return nil
	}
	if d.NodeType == NodeCache {
		dir = fmt.Sprintf("/var/lib/drone/cache/%s:/cache", dir)
		d.Volumes = []string{dir}
	}
	return nil
}

func CacheFunc(dir string) RuleFunc {
	return func(n Node) error {
		return Cache(n, dir)
	}
}

// expandImage expands an alias plugin name to use a
// fully qualified image name.
func expandImage(image string) string {
	if !strings.Contains(image, "/") {
		image = path.Join("plugins", "drone-"+image)
	}
	return strings.Replace(image, "_", "-", -1)
}

// expandImageDefault returns the default image if none
// is specified in the Yaml. If an image is specified,
// it expands the alias.
func expandImageDefault(image, defaultImage string) string {
	if len(image) == 0 {
		return defaultImage
	}
	return expandImage(image)
}

// expandImageTag is a helper function that automatically
// expands the image to include the :latest tag if not present
func expandImageTag(image string) string {
	if strings.Contains(image, "@") {
		return image
	}
	n := strings.LastIndex(image, ":")
	if n < 0 {
		return image + ":latest"
	}
	if tag := image[n+1:]; strings.Contains(tag, "/") {
		return image + ":latest"
	}
	return image
}
