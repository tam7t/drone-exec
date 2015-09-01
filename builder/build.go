package builder

import (
	"errors"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/drone/drone-exec/builder/parse"
	"github.com/drone/drone-exec/docker"
	"github.com/samalba/dockerclient"
)

var ErrNoImage = errors.New("Yaml must specify an image for every step")

// Default clone plugin.
const DefaultCloner = "plugins/drone-git"

// Default cache plugin.
const DefaultCacher = "plugins/drone-cache"

type Build struct {
	tree  *parse.Tree
	flags parse.NodeType

	parseFuncs []func()
	execFuncs  []func()
}

func (b *Build) Run(state *State) error {
	return b.RunNode(state, 0)
}

func (b *Build) RunNode(state *State, flags parse.NodeType) error {
	b.flags = flags
	return b.walk(b.tree.Root, state)
}

func (b *Build) walk(node parse.Node, state *State) (err error) {

	switch node := node.(type) {
	case *parse.ListNode:
		for _, node := range node.Nodes {
			err = b.walk(node, state)
			if err != nil {
				break
			}
		}

	case *parse.FilterNode:
		if isMatch(node, state) {
			b.walk(node.Node, state)
		} else {
			log.Warnf("skipping step. filter conditions not met")
		}

	case *parse.DockerNode:

		if shouldSkip(b.flags, node.NodeType) {
			log.Warnf("skipping step. flag to run build not exists")
			break
		}

		switch node.Type() {

		case parse.NodeBuild:
			// run setup
			node.Vargs = map[string]interface{}{}
			node.Vargs["commands"] = node.Commands

			conf := toContainerConfig(node)
			conf.Cmd = toCommand(state, node)
			conf.Image = "plugins/drone-build"
			info, err := docker.Run(state.Client, conf, node.Pull)
			if err != nil {
				state.Exit(255)
			} else if info.State.ExitCode != 0 {
				state.Exit(info.State.ExitCode)
			}

			// run build
			conf = toContainerConfig(node)
			conf.Entrypoint = []string{"/bin/sh", "-e"}
			conf.Cmd = []string{"/drone/bin/build.sh"}
			info, err = docker.Run(state.Client, conf, node.Pull)
			if err != nil {
				state.Exit(255)
			} else if info.State.ExitCode != 0 {
				state.Exit(info.State.ExitCode)
			}

		case parse.NodeCompose:
			conf := toContainerConfig(node)
			_, err := docker.Start(state.Client, conf, node.Pull)
			if err != nil {
				state.Exit(255)
			}

		default:
			conf := toContainerConfig(node)
			conf.Cmd = toCommand(state, node)
			info, err := docker.Run(state.Client, conf, node.Pull)
			if err != nil {
				state.Exit(255)
			} else if info.State.ExitCode != 0 {
				state.Exit(info.State.ExitCode)
			}
		}
	}

	return nil
}

func expectMatch() {

}

func maybeResolveImage() {}

func maybeEscalate(conf dockerclient.ContainerConfig, node *parse.DockerNode) {
	if node.Image == "plugins/drone-docker" {
		return
	}
	conf.Volumes = nil
	conf.HostConfig.NetworkMode = ""
	conf.HostConfig.Privileged = true
	conf.Entrypoint = []string{}
	conf.Cmd = []string{}
}

// shouldSkip is a helper function that returns true if
// node execution should be skipped. This happens when
// the build is executed for a subset of build steps.
func shouldSkip(flags parse.NodeType, nodeType parse.NodeType) bool {
	return flags != 0 && flags&nodeType == 0
}

// shouldEscalate is a helper function that returns true
// if the plugin should be escalated to start the container
// in privileged mode.
func shouldEscalate(node *parse.DockerNode) bool {
	return node.Image == "plugins/drone-docker"
}

// resolveImage is a helper function that resolves the docker
// image name. Plugins may use a short, alias name which needs
// to be expanded.
func resolveImage(node *parse.DockerNode) error {
	switch node.NodeType {
	case parse.NodeBuild, parse.NodeCompose:
		break
	case parse.NodeClone:
		node.Image = expandImageDefault(node.Image, DefaultCloner)
	case parse.NodeCache:
		node.Image = expandImageDefault(node.Image, DefaultCacher)
	default:
		node.Image = expandImage(node.Image)
	}
	if len(node.Image) == 0 {
		return ErrNoImage
	}
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
