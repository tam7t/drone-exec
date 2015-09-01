package parse

// import (
// 	"errors"
// 	"path"
// 	"strings"
// )

// var ErrNoImage = errors.New("Yaml must specify an image for every step")

// // Default clone plugin.
// const DefaultCloner = "plugins/drone-git"

// // Default cache plugin.
// const DefaultCacher = "plugins/drone-cache"

// // ParseFunc is a function used to validate or modify
// // the yaml during the parsing process. This is helpful
// // for tasks like cleansing the yaml, removing restricted
// // configuration options, etc.
// type ParseFunc func(Node) error

// // Default parser functions.
// var DefaultFuncs = []ParseFunc{resolveImage}

// // resolveImage is a parser function that resolves the docker
// // node image. Plugins may use a short, alias name which needs
// // to be expanded.
// func resolveImage(n Node) error {
// 	d, ok := n.(*DockerNode)
// 	if !ok {
// 		return nil
// 	}
// 	switch d.NodeType {
// 	case NodeBuild, NodeCompose:
// 		break
// 	case NodeClone:
// 		d.Image = expandImageDefault(d.Image, DefaultCloner)
// 	case NodeCache:
// 		d.Image = expandImageDefault(d.Image, DefaultCacher)
// 	default:
// 		d.Image = expandImage(d.Image)
// 	}
// 	if len(d.Image) == 0 {
// 		return ErrNoImage
// 	}
// 	return nil
// }

// // expandImage expands an alias plugin name to use a
// // fully qualified image name.
// func expandImage(image string) string {
// 	if !strings.Contains(image, "/") {
// 		image = path.Join("plugins", "drone-"+image)
// 	}
// 	return strings.Replace(image, "_", "-", -1)
// }

// // expandImageDefault returns the default image if none
// // is specified in the Yaml. If an image is specified,
// // it expands the alias.
// func expandImageDefault(image, defaultImage string) string {
// 	if len(image) == 0 {
// 		return defaultImage
// 	}
// 	return expandImage(image)
// }
