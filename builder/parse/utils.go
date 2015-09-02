package parse

import (
	"errors"
	"path"
	"strings"
)

var ErrNoImage = errors.New("Yaml must specify an image for every step")

// Default clone plugin.
const DefaultCloner = "plugins/drone-git"

// Default cache plugin.
const DefaultCacher = "plugins/drone-cache"

// resolveImage is a parser function that resolves the docker
// node image. Plugins may use a short, alias name which needs
// to be expanded.
func resolveImage(d *DockerNode) error {
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
		return ErrNoImage
	}
	d.Image = expandImageTag(d.Image)
	return nil
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
