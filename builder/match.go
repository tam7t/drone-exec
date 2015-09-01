package builder

import (
	"path"
	"strings"

	"github.com/drone/drone-exec/builder/parse"
)

// isMatch is a helper function that returns true if
// all criteria is matched.
func isMatch(node *parse.FilterNode, s *State) bool {
	return matchBranch(node.Branch, s.Build.Commit.Branch) &&
		matchMatrix(node.Matrix, s.Job.Environment) &&
		matchRepo(node.Repo, s.Repo.FullName)
}

// matchBranch is a helper function that returns true
// if all_branches is true. Else it returns false if a
// branch condition is specified, and the branch does
// not match.
func matchBranch(want, got string) bool {
	if len(want) == 0 {
		return true
	}
	if strings.HasPrefix(got, "refs/heads/") {
		got = got[11:]
	}
	match, _ := path.Match(want, got)
	return match
}

// matchRepo is a helper function that returns false
// if this task is only intended for a named repo,
// the current build does not match that repo.
//
// This is useful when you want to prevent forks from
// executing deployment, publish or notification steps.
func matchRepo(want, got string) bool {
	if len(want) == 0 {
		return true
	}
	return got == want
}

// matchMatrix is a helper function that returns false
// to limit steps to only certain matrix axis.
func matchMatrix(want, got map[string]string) bool {
	if len(want) == 0 {
		return true
	}
	for k, v := range want {
		if got[k] != v {
			return false
		}
	}
	return true
}
