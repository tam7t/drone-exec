package builder

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/drone/drone-exec/builder/parse"
	"github.com/drone/drone-plugin-go/plugin"
	"github.com/samalba/dockerclient"
)

// helper function that converts the build step to
// a containerConfig for use with the dockerclient
func toContainerConfig(n *parse.DockerNode) *dockerclient.ContainerConfig {
	config := &dockerclient.ContainerConfig{
		Image:      n.Image,
		Env:        n.Environment,
		Cmd:        n.Command,
		Entrypoint: n.Entrypoint,
		HostConfig: dockerclient.HostConfig{
			Privileged:  n.Privileged,
			NetworkMode: n.Net,
		},
	}

	if len(config.Entrypoint) == 0 {
		config.Entrypoint = nil
	}

	config.Volumes = map[string]struct{}{}
	for _, path := range n.Volumes {
		if strings.Index(path, ":") == -1 {
			continue
		}
		parts := strings.Split(path, ":")
		config.Volumes[parts[1]] = struct{}{}
		config.HostConfig.Binds = append(config.HostConfig.Binds, path)
	}

	return config
}

// helper function to inject drone-specific environment
// variables into the container.
func toEnv(s *State) map[string]string {
	return map[string]string{
		"CI":           "true",
		"BUILD_DIR":    s.Workspace.Path,
		"BUILD_ID":     strconv.Itoa(s.Build.Number),
		"BUILD_NUMBER": strconv.Itoa(s.Build.Number),
		"JOB_NAME":     s.Repo.FullName,
		"WORKSPACE":    s.Workspace.Path,
		"GIT_BRANCH":   s.Build.Commit.Branch,
		"GIT_COMMIT":   s.Build.Commit.Sha,

		"DRONE":        "true",
		"DRONE_REPO":   s.Repo.FullName,
		"DRONE_BUILD":  strconv.Itoa(s.Build.Number),
		"DRONE_BRANCH": s.Build.Commit.Branch,
		"DRONE_COMMIT": s.Build.Commit.Sha,
		"DRONE_DIR":    s.Workspace.Path,
	}
}

// helper function to encode the build step to
// a json string. Primarily used for plugins, which
// expect a json encoded string in stdin or arg[1].
func toCommand(s *State, n *parse.DockerNode) []string {
	p := payload{
		Workspace: s.Workspace,
		Repo:      s.Repo,
		Build:     s.Build,
		Job:       s.Job,
		Vargs:     n.Vargs,

		Clone: &plugin.Clone{
			Origin: s.Repo.Clone,
			Remote: s.Repo.Clone,
			Branch: s.Build.Commit.Branch,
			Sha:    s.Build.Commit.Sha,
			Ref:    s.Build.Commit.Ref,
			Dir:    s.Workspace.Path,
		},
	}
	p.System = &plugin.System{
		Version: s.System.Version,
		Link:    s.System.Link,
	}
	b, _ := json.Marshal(p)
	return []string{string(b)}
}

// payload represents the payload of a plugin
// that is serialized and sent to the plugin in JSON
// format via stdin or arg[1].
type payload struct {
	Workspace *plugin.Workspace `json:"workspace"`
	System    *plugin.System    `json:"system"`
	Repo      *plugin.Repo      `json:"repo"`
	Build     *plugin.Build     `json:"build"`
	Job       *plugin.Job       `json:"job"`
	Clone     *plugin.Clone     `json:"clone"`

	Vargs map[string]interface{} `json:"vargs"`
}
