package builder

import (
	"encoding/json"
	"fmt"
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
func toEnv(s *State) []string {
	var envs []string

	envs = append(envs, "CI=true")
	envs = append(envs, "DRONE=true")
	envs = append(envs, "CI_NAME=drone")

	envs = append(envs, fmt.Sprintf("DRONE_DIR=%s", s.Workspace.Path))
	envs = append(envs, fmt.Sprintf("DRONE_REPO=%s", s.Repo.FullName))
	envs = append(envs, fmt.Sprintf("CI_REPO=%s", s.Repo.FullName))

	// environment variables specific to the job
	envs = append(envs, fmt.Sprintf("DRONE_JOB_ID=%d", s.Job.ID))
	envs = append(envs, fmt.Sprintf("DRONE_JOB_NUMBER=%d", s.Job.Number))
	envs = append(envs, fmt.Sprintf("CI_JOB_ID=%d", s.Job.ID))
	envs = append(envs, fmt.Sprintf("CI_JOB_NUMBER=%d", s.Job.Number))

	// environment variables specific to the build
	envs = append(envs, fmt.Sprintf("DRONE_BUILD_NUMBER=%d", s.Build.Number))
	envs = append(envs, fmt.Sprintf("DRONE_BUILD_DIR=%s", s.Workspace.Path))
	envs = append(envs, fmt.Sprintf("DRONE_BRANCH=%s", s.Build.Commit.Branch))
	envs = append(envs, fmt.Sprintf("DRONE_COMMIT=%s", s.Build.Commit.Sha))
	envs = append(envs, fmt.Sprintf("CI_BRANCH=%s", s.Build.Commit.Branch))
	envs = append(envs, fmt.Sprintf("CI_BUILD_DIR=%s", s.Workspace.Path))
	envs = append(envs, fmt.Sprintf("CI_BUILD_NUMBER=%d", s.Build.Number))
	envs = append(envs, fmt.Sprintf("CI_COMMIT=%s", s.Build.Commit.Sha))

	envs = append(envs, fmt.Sprintf("CI_BUILD_URL=%s/%s/%d", s.System.Link, s.Repo.FullName, s.Build.Number))

	// environment variables specific to the pull request
	if s.Build.PullRequest != nil {
		envs = append(envs, fmt.Sprintf("CI_PULL_REQUEST=%d", s.Build.PullRequest.Number))
		envs = append(envs, fmt.Sprintf("DRONE_PULL_REQUEST=%d", s.Build.PullRequest.Number))
	}

	// environment variables for the current matrix axis
	for key, val := range s.Job.Environment {
		envs = append(envs, fmt.Sprintf("%s=%s", key, val))
	}

	return envs
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
	}
	p.System = &plugin.System{
		Version: s.System.Version,
		Link:    s.System.Link,
	}
	b, _ := json.Marshal(p)
	return []string{"--", string(b)}
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

	Vargs map[string]interface{} `json:"vargs"`
}
