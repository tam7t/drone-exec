package runner

import (
	"fmt"
	"path"
	"strings"

	"github.com/drone/drone-exec/parser"
	"github.com/drone/drone-plugin-go/plugin"
)

// isMatch is a helper function that returns true if
// all criteria is matched.
func isMatch(node *parser.FilterNode, s *State) bool {

	// TODO this code should be stored inside the
	// build object itself. So we need to add
	// to the database
	event := "push"
	if s.Build.PullRequest != nil {
		event = "pull_request"
	} else if strings.HasPrefix(s.Build.Commit.Ref, "refs/tags") {
		event = "tag"
	}

	return matchBranch(node.Branch, s.Build.Commit.Branch) &&
		matchMatrix(node.Matrix, s.Job.Environment) &&
		matchRepo(node.Repo, s.Repo.FullName) &&
		matchSuccess(node.Success, s.Build.Status) &&
		matchFailure(node.Success, s.Build.Status) &&
		matchEvent(node.Event, event)
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

// matchEvent is a helper function that returns false
// if this task is only intended for a specific repository
// event not matched by the current build. For example,
// only executing a build for `tags` or `pull_requests`
func matchEvent(want, got string) bool {
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

func matchSuccess(toggle, status string) bool {
	ok, err := parseBool(toggle)
	if err != nil {
		return true
	}
	return ok && status == plugin.StateSuccess
}

func matchFailure(toggle, status string) bool {
	ok, err := parseBool(toggle)
	if err != nil {
		return true
	}
	return ok && status != plugin.StateSuccess
}

func parseBool(str string) (value bool, err error) {
	switch str {
	case "true", "TRUE", "True", "On", "ON", "on":
		return true, nil
	case "false", "FALSE", "False", "Off", "off", "OFF":
		return false, nil
	}
	return false, fmt.Errorf("Error parsing boolean %s", str)
}
