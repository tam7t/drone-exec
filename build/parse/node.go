package parse

import "github.com/drone/drone-exec/yaml"

// NodeType identifies the type of a parse tree node.
type NodeType uint

// Type returns itself and provides an easy default
// implementation for embedding in a Node. Embedded
// in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeList NodeType = 1 << iota
	NodeFilter
	NodeBuild
	NodeCache
	NodeClone
	NodeDeploy
	NodeCompose
	NodeNotify
	NodePublish
)

// Nodes.

type Node interface {
	Type() NodeType
}

// ListNode holds a sequence of nodes.
type ListNode struct {
	NodeType
	Nodes []Node // nodes executed in lexical order.
}

// Append appends a node to the list.
func (l *ListNode) append(n ...Node) {
	l.Nodes = append(l.Nodes, n...)
}

func newListNode() *ListNode {
	return &ListNode{NodeType: NodeList}
}

// DockerNode represents a Docker container that
// should be laucned as part of the build process.
type DockerNode struct {
	NodeType

	Image       string
	Pull        bool
	Privileged  bool
	Environment []string
	Entrypoint  []string
	Command     []string
	Commands    []string
	Volumes     []string
	Net         string
	Vargs       map[string]interface{}
}

func newDockerNode(typ NodeType, c yaml.Container) (*DockerNode, error) {
	var node = DockerNode{
		NodeType:    typ,
		Image:       c.Image,
		Pull:        c.Pull,
		Privileged:  c.Privileged,
		Environment: c.Environment.Slice(),
		Entrypoint:  c.Entrypoint.Slice(),
		Command:     c.Command.Slice(),
		Volumes:     c.Volumes,
		Net:         c.Net,
	}
	var err = resolveImage(&node)
	return &node, err
}

func newPluginNode(typ NodeType, p yaml.Plugin) (*DockerNode, error) {
	node, err := newDockerNode(typ, p.Container)
	node.Vargs = p.Vargs
	return node, err
}

func newBuildNode(typ NodeType, b yaml.Build) (*DockerNode, error) {
	node, err := newDockerNode(typ, b.Container)
	node.Commands = b.Commands
	return node, err
}

// FilterNode represents a conditional step used to
// filter nodes. If conditions are met the child
// node is executed.
type FilterNode struct {
	NodeType

	Repo    string
	Branch  string
	Success string
	Failure string
	Matrix  map[string]string

	Node Node // Node to execution if conditions met
}

func newFilterNode(p yaml.Plugin) *FilterNode {
	return &FilterNode{
		NodeType: NodeFilter,
		Repo:     p.Filter.Repo,
		Branch:   p.Filter.Branch,
		Matrix:   p.Filter.Matrix,
		Success:  p.Filter.Success,
		Failure:  p.Filter.Failure,
	}
}
