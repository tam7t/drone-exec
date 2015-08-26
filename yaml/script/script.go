package script

import (
	"bytes"
	"io"

	"github.com/drone/drone-exec/yaml"
)

// TransformFunc represents a transform function used when
// parsing the yaml file and setting default values and state.
type TransformFunc func(conf *yaml.Config)

// WriterFunc represents a transform function used to
// transform the the yaml file into a build script.
type WriterFunc func(io.Writer)

// Transform returns a yaml transform function responsible for
// executing a set of build transforms, altering the Build step
// to run the build script in the entrypoint.
func Transform(funcs []WriterFunc) TransformFunc {
	return func(conf *yaml.Config) {
		var buf bytes.Buffer
		WriteAll(&buf, funcs)
		conf.Build.Entrypoint = []string{"/bin/sh"}
		conf.Build.Command = []string{"-e", "-c", wrapCommand(buf.Bytes())}
	}
}

// WriteAll executes all writer transforms.
func WriteAll(w io.Writer, funcs []WriterFunc) {
	for _, fn := range funcs {
		fn(w)
	}
}

// WriteSetup returns a transform to add setup commands
// to the build script.
func WriteSetup() WriterFunc {
	return func(w io.Writer) {
		io.WriteString(w, setupScript)
	}
}

// WriteTeardown returns a transform to add teardown
// commands to the build script.
func WriteTeardown() WriterFunc {
	return func(w io.Writer) {
		io.WriteString(w, teardownScript)
	}
}

// WriteNetrc returns a transform to add commands to the
// build script that setup a .netrc file.
func WriteNetrc(machine, login, password string) WriterFunc {
	return func(w io.Writer) {
		io.WriteString(w, writeNetrc(machine, login, password))
	}
}

// WriteKey returns a transform to add commands to the
// build script that setup an id_rsa file and configuration.
func WriteKey(key []byte) WriterFunc {
	return func(w io.Writer) {
		io.WriteString(w, writeKey(key))
	}
}

// WriteCmds returns a transform to add the user-defined
// commands to the build script, with ability to echo
// those commands to the build output prior to execution.
func WriteCmds(conf *yaml.Config) WriterFunc {
	return func(w io.Writer) {
		if conf.Build == nil || conf.Build.Config == nil {
			return
		}
		cmdv, ok := conf.Build.Config["commands"]
		if !ok {
			return
		}
		cmds, ok := cmdv.([]string)
		if !ok {
			return
		}
		for _, cmd := range cmds {
			io.WriteString(w, traceCommand(cmd))
		}
	}
}
