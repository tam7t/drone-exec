package parse

import "github.com/drone/drone-exec/yaml"

// Tree is the representation of a parsed build
// configuraiton Yaml file.
type Tree struct {
	Root *ListNode
}

// New allocates a new parse tree.
func New() *Tree {
	return &Tree{
		Root: &ListNode{NodeType: NodeList},
	}
}

// Parse parses the Yaml build definition file
// and returns an execution Tree.
func Parse(raw string) (*Tree, error) {
	conf, err := yaml.ParseString(raw)
	if err != nil {
		return nil, err
	}
	return Load(conf)
}

// Load loads the Yaml build definition structure
// and returns an execution Tree.
func Load(conf *yaml.Config) (*Tree, error) {
	var tree = New()
	var err error

	// Cache.
	err = appendCache(tree.Root, conf.Cache)
	if err != nil {
		return nil, err
	}

	// Clone.
	err = appendPlugin(tree.Root, NodeClone, conf.Clone)
	if err != nil {
		return nil, err
	}

	// Compose.
	err = appendCompose(tree.Root, conf.Compose.Slice())
	if err != nil {
		return nil, err
	}

	// Build
	err = appendBuild(tree.Root, conf.Build)
	if err != nil {
		return nil, err
	}

	// Publish.
	appendPlugin(tree.Root, NodePublish, conf.Publish.Slice()...)
	if err != nil {
		return nil, err
	}

	// Deploy.
	appendPlugin(tree.Root, NodeDeploy, conf.Deploy.Slice()...)
	if err != nil {
		return nil, err
	}

	// Plugin.
	appendPlugin(tree.Root, NodeNotify, conf.Notify.Slice()...)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func appendPlugin(list *ListNode, typ NodeType, plugins ...yaml.Plugin) error {
	for _, plugin := range plugins {
		node, err := newPluginNode(typ, plugin)
		if err != nil {
			return err
		}
		fnode := newFilterNode(plugin)
		fnode.Node = node
		list.append(fnode)
	}
	return nil
}

func appendBuild(list *ListNode, build yaml.Build) error {
	node, err := newBuildNode(NodeBuild, build)
	if err != nil {
		return err
	}
	list.append(node)
	return nil
}

func appendCache(list *ListNode, cache yaml.Plugin) error {
	if len(cache.Vargs) == 0 {
		return nil
	}
	return appendPlugin(list, NodeCache, cache)
}

func appendCompose(list *ListNode, plugins []yaml.Container) error {
	for _, plugin := range plugins {
		node, err := newDockerNode(NodeCompose, plugin)
		if err != nil {
			return err
		}
		list.append(node)
	}
	return nil
}

// func (t *Tree) validate(node Node) error {
// 	switch node := node.(type) {
// 	case *ListNode:
// 		for _, node := range node.Nodes {
// 			err := t.walk(node)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	case *FilterNode:
// 		return t.walk(node.Node)
// 	default:
// 		return t.validate(node)
// 	}
// 	return nil
// }

// func (t *Tree) validate(n Node) error {
//  // run validation functions
// 	return nil
// }
