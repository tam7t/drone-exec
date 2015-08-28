package script

import (
	"bytes"
	"io"

	"github.com/drone/drone-exec/yaml"
	"github.com/drone/drone-exec/yaml/rules"
)

// WriterFunc represents a transform function used to
// transform the the yaml file into a build script.
type WriterFunc func(io.Writer)

// RuleFunc returns a yaml rule function responsible for
// executing a set of build transforms, altering the Build step
// to run the build script in the entrypoint.
func RuleFunc(funcs []WriterFunc) rules.RuleFunc {
	return func(conf *yaml.Config) error {
		var buf bytes.Buffer
		WriteAll(&buf, funcs)
		conf.Build.Entrypoint = []string{"/bin/sh", "-e", "-c"}
		conf.Build.Command = []string{wrapCommand(buf.Bytes())}
		return nil
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
