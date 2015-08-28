package yaml

import (
	"path/filepath"
	"strings"
)

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
	Repo   string // Indicates the step should run only for this repo (useful for forks)
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

// MatchRepo is a helper function that returns false
// if this task is only intended for a named repo,
// the current build does not match that repo.
//
// This is useful when you want to prevent forks from
// executing deployment, publish or notification steps.
func (c *When) MatchRepo(repo string) bool {
	if len(c.Repo) == 0 {
		return true
	}
	return repo == c.Repo
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
