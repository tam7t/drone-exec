package parser

import "github.com/drone/drone-exec/yaml"

// Tree is the representation of a parsed build
// configuraiton Yaml file.
type Tree struct {
	Root  *ListNode
	rules []RuleFunc
}

// newTree allocates a new parse tree.
func newTree(rules []RuleFunc) *Tree {
	return &Tree{
		Root:  &ListNode{NodeType: NodeList},
		rules: rules,
	}
}

// Parse parses the Yaml build definition file
// and returns an execution Tree.
func Parse(raw string, rules []RuleFunc) (*Tree, error) {
	conf, err := yaml.ParseString(raw)
	if err != nil {
		return nil, err
	}
	return Load(conf, rules)
}

// Load loads the Yaml build definition structure
// and returns an execution Tree.
func Load(conf *yaml.Config, rules []RuleFunc) (*Tree, error) {
	var tree = newTree(rules)
	var err error

	// Cache.
	err = tree.appendCache(conf.Cache)
	if err != nil {
		return nil, err
	}

	// Clone.
	err = tree.appendPlugin(NodeClone, conf.Clone)
	if err != nil {
		return nil, err
	}

	// Compose.
	err = tree.appendCompose(conf.Compose.Slice())
	if err != nil {
		return nil, err
	}

	// Build
	err = tree.appendBuild(conf.Build.Slice())
	if err != nil {
		return nil, err
	}

	// Publish.
	err = tree.appendPlugin(NodePublish, conf.Publish.Slice()...)
	if err != nil {
		return nil, err
	}

	// Deploy.
	err = tree.appendPlugin(NodeDeploy, conf.Deploy.Slice()...)
	if err != nil {
		return nil, err
	}

	// Plugin.
	err = tree.appendPlugin(NodeNotify, conf.Notify.Slice()...)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func (t *Tree) appendPlugin(typ NodeType, plugins ...yaml.Plugin) error {
	for _, plugin := range plugins {
		node := newPluginNode(typ, plugin)
		for _, rule := range t.rules {
			err := rule(node)
			if err != nil {
				return err
			}
		}
		fnode := newFilterNode(plugin.Filter)
		fnode.Node = node
		// TODO: we should apply rules to all nodes in
		// the tree AFTER the entire tree is constructed.
		for _, rule := range t.rules {
			err := rule(fnode)
			if err != nil {
				return err
			}
		}
		t.Root.append(fnode)
	}
	return nil
}

func (t *Tree) appendBuild(builds []yaml.Build) error {
	for _, build := range builds {
		node := newBuildNode(NodeBuild, build)
		for _, rule := range t.rules {
			if err := rule(node); err != nil {
				return err
			}
		}

		fnode := newFilterNode(build.Filter)
		fnode.Node = node
		for _, rule := range t.rules {
			if err := rule(fnode); err != nil {
				return err
			}
		}
		t.Root.append(fnode)
	}
	return nil
}

func (t *Tree) appendCache(cache yaml.Plugin) error {
	if len(cache.Vargs) == 0 {
		return nil
	}
	return t.appendPlugin(NodeCache, cache)
}

func (t *Tree) appendCompose(plugins []yaml.Container) error {
	for _, plugin := range plugins {
		node := newDockerNode(NodeCompose, plugin)
		for _, rule := range t.rules {
			err := rule(node)
			if err != nil {
				return err
			}
		}
		t.Root.append(node)
	}
	return nil
}
