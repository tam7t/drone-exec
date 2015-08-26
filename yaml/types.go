package yaml

import (
	"path/filepath"
	"strings"
)

// Config represents a repository build configuration.
type Config struct {
	Setup *Step
	Clone *Step
	Build *Step

	Compose map[string]*Step
	Publish map[string]*Step
	Deploy  map[string]*Step
	Notify  map[string]*Step

	Workspace Workspace
}

// Step represents a step in the build process, including
// the execution environment and parameters.
type Step struct {
	Image       string
	Pull        bool
	Privileged  bool
	Environment []string
	Entrypoint  []string
	Command     []string
	Volumes     []string
	Cache       []string
	WorkingDir  string `yaml:"working_dir"`
	NetworkMode string `yaml:"net"`

	// Condition represents a set of conditions that must
	// be met in order to execute this step.
	When *When

	// Config represents the unique configuration details
	// for each plugin.
	Config map[string]interface{} `yaml:"config,inline"`
}

// When represents a set of conditions that must
// be met in order to proceed with a build or build step.
type When struct {
	Owner  string // Indicates the step should run only for this repo (useful for forks)
	Branch string // Indicates the step should run only for this branch

	// Indicates the step should only run when the following
	// matrix values are present for the sub-build.
	Matrix map[string]string
}

// MatchBranch is a helper function that returns true
// if all_branches is true. Else it returns false if a
// branch condition is specified, and the branch does
// not match.
func (c *When) MatchBranch(branch string) bool {
	if len(c.Branch) == 0 {
		return true
	}
	if strings.HasPrefix(branch, "refs/heads/") {
		branch = branch[11:]
	}
	match, _ := filepath.Match(c.Branch, branch)
	return match
}

// MatchOwner is a helper function that returns false
// if an owner condition is specified and the repository
// owner does not match.
//
// This is useful when you want to prevent forks from
// executing deployment, publish or notification steps.
func (c *When) MatchOwner(owner string) bool {
	if len(c.Owner) == 0 {
		return true
	}
	parts := strings.Split(owner, "/")
	switch len(parts) {
	case 2:
		return c.Owner == parts[0]
	case 3:
		return c.Owner == parts[1]
	default:
		return c.Owner == owner
	}
}

// MatchMatrix is a helper function that returns false
// to limit steps to only certain matrix axis.
func (c *When) MatchMatrix(matrix map[string]string) bool {
	if len(c.Matrix) == 0 {
		return true
	}
	for k, v := range c.Matrix {
		if matrix[k] != v {
			return false
		}
	}
	return true
}

// Workspace defines the build's workspace inside the
// container. This helps the plugin understand locate
// the source code directory.
type Workspace struct {
	Path  string
	Cache []string
}
