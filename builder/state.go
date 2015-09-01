package builder

import (
	"io"
	"sync"

	"github.com/drone/drone-plugin-go/plugin"
	"github.com/samalba/dockerclient"
)

type State struct {
	sync.Mutex

	Repo      *plugin.Repo
	Build     *plugin.Build
	Job       *plugin.Job
	System    *plugin.System
	Workspace *plugin.Workspace

	// Client is an instance of the Docker client
	// used to spawn container tasks.
	Client dockerclient.Client

	Stdout, Stderr io.Writer
}

// Exit writes the function as having failed but
// continues execution.
func (s *State) Exit(code int) {
	s.Lock()
	defer s.Unlock()

	if code != 0 { // never override non-zero exit
		s.Job.ExitCode = code
	}
}

// ExitCode reports the build exit code. A non-zero
// value indicates the build exited with errors.
func (s *State) ExitCode() int {
	s.Lock()
	defer s.Unlock()

	return s.Job.ExitCode
}

// IsFailed returns true if the build has a
// non-zero exit code.
func (s *State) IsFailed() bool {
	return s.ExitCode() != 0
}
