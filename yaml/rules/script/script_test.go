package script

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/drone/drone-exec/yaml"
	"github.com/franela/goblin"
)

func Test_Writer(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Script writer", func() {

		g.It("Should add setup script", func() {
			var buf bytes.Buffer
			fn := WriteSetup()
			fn(&buf)

			g.Assert(buf.String()).Equal(setupScript)
		})

		g.It("Should add teardown script", func() {
			var buf bytes.Buffer
			fn := WriteTeardown()
			fn(&buf)

			g.Assert(buf.String()).Equal(teardownScript)
		})

		g.It("Should add netrc script", func() {
			var buf bytes.Buffer
			fn := WriteNetrc("foo", "bar", "baz")
			fn(&buf)

			g.Assert(buf.String()).Equal(netrc)
		})

		g.It("Should not add netrc script if empty machine", func() {
			var buf bytes.Buffer
			fn := WriteNetrc("", "bar", "baz")
			fn(&buf)

			g.Assert(buf.String()).Equal("")
		})

		g.It("Should add rsa key script", func() {
			var buf bytes.Buffer
			fn := WriteKey([]byte("-----BEGIN RSA PRIVATE KEY----- MIIEpQIBAAKCAQEA3Tz2..."))
			fn(&buf)
			g.Assert(buf.String()).Equal(keys)
		})

		g.It("Should not add rsa key script if empty key", func() {
			var buf bytes.Buffer
			fn := WriteKey([]byte(""))
			fn(&buf)

			g.Assert(buf.String()).Equal("")
		})

		g.It("Should write commands with trace", func() {
			conf := yaml.Config{
				Build: &yaml.Step{
					Config: map[string]interface{}{
						"commands": []string{"echo hello", "echo goodbye"},
					},
				},
			}

			var buf bytes.Buffer
			fn := WriteCmds(&conf)
			fn(&buf)

			g.Assert(strings.Trim(buf.String(), "\n")).Equal(strings.Trim(cmds, "\n"))
		})

		g.It("Should anticipate a nil Build when writing commands", func() {
			var buf bytes.Buffer
			fn := WriteCmds(&yaml.Config{})
			fn(&buf)

			g.Assert(buf.String()).Equal("")
		})

		g.It("Should anticipate a nil Build Config when writing commands", func() {
			var buf bytes.Buffer
			fn := WriteCmds(&yaml.Config{Build: &yaml.Step{}})
			fn(&buf)

			g.Assert(buf.String()).Equal("")
		})

		g.It("Should anticipate a nil Build Config when writing commands", func() {
			var buf bytes.Buffer
			build := &yaml.Step{}
			build.Config = map[string]interface{}{}
			fn := WriteCmds(&yaml.Config{Build: build})
			fn(&buf)

			g.Assert(buf.String()).Equal("")
		})

		g.It("Should anticipate a malformed Build Config when writing commands", func() {
			var buf bytes.Buffer
			build := &yaml.Step{}
			build.Config = map[string]interface{}{"commands": 0}
			fn := WriteCmds(&yaml.Config{Build: build})
			fn(&buf)

			g.Assert(buf.String()).Equal("")
		})

		g.It("Should process multiple writers", func() {
			var buf bytes.Buffer
			WriteAll(&buf, []WriterFunc{WriteSetup()})

			g.Assert(buf.String()).Equal(setupScript)
		})
	})
}

func Test_Rule(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Script transform", func() {

		conf := &yaml.Config{}
		conf.Build = &yaml.Step{}

		var funcs []WriterFunc
		funcs = append(funcs, WriteSetup())

		fn := RuleFunc(funcs)
		fn(conf)

		g.It("Should modify entrypoint", func() {
			g.Assert(conf.Build.Entrypoint).Equal([]string{"/bin/sh", "-e", "-c"})
		})

		g.It("Should modify embed base64 command", func() {
			cmd := base64.StdEncoding.EncodeToString([]byte(setupScript))
			cmd = fmt.Sprintf("echo %s | base64 -d | $SHELL", cmd)
			g.Assert(conf.Build.Command).Equal([]string{cmd})
		})
	})
}

var netrc = `
cat <<EOF > $HOME/.netrc
machine foo
login bar
password baz
EOF
chmod 0600 $HOME/.netrc
`

var keys = `
mkdir -p $HOME/.ssh
cat <<EOF > $HOME/.ssh/id_rsa
-----BEGIN RSA PRIVATE KEY----- MIIEpQIBAAKCAQEA3Tz2...
EOF
chmod 0700 $HOME/.ssh

cat <<EOF > $HOME/.ssh/config
StrictHostKeyChecking no
EOF

mkdir -p /etc/apt/apt.conf.d
cat <<EOF > /etc/apt/apt.conf.d/90forceyes
APT::Get::Assume-Yes "true";APT::Get::force-yes "true";
EOF
`

var cmds = `
echo JCBlY2hvIGhlbGxvCg== | base64 -d
echo hello
echo JCBlY2hvIGdvb2RieWUK | base64 -d
echo goodbye
`
