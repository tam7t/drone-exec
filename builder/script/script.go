package script

import (
	"bytes"

	"github.com/drone/drone-exec/builder/parse"
	"github.com/drone/drone-plugin-go/plugin"
	"github.com/samalba/dockerclient"
)

// Encode encodes the build script as a command in the
// provided Container config. For linux, the build script
// is embedded as the container entrypoint command, base64
// encoded as a one-line script.
func Encode(w *plugin.Workspace, c *dockerclient.ContainerConfig, n *parse.DockerNode) {
	var buf bytes.Buffer
	buf.WriteString(setupScript)

	if w != nil && w.Keys != nil && w.Netrc != nil {
		buf.WriteString(writeKey(
			w.Keys.Private,
		))
		buf.WriteString(writeNetrc(
			w.Netrc.Machine,
			w.Netrc.Login,
			w.Netrc.Password,
		))
	}

	buf.WriteString(writeCmds(n.Commands))
	buf.WriteString(teardownScript)

	c.Entrypoint = entrypoint
	c.Cmd = []string{encode(buf.Bytes())}
}
