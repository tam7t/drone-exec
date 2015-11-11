package inject

import (
	"testing"

	"github.com/franela/goblin"
	"gopkg.in/yaml.v2"
)

func Test_Inject(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Inject params", func() {

		g.It("Should replace vars with $$", func() {
			s := "echo $$FOO $BAR"
			m := map[string]string{}
			m["FOO"] = "BAZ"
			g.Assert("echo BAZ $BAR").Equal(Inject(s, m))
		})

		g.It("Should not replace vars with single $", func() {
			s := "echo $FOO $BAR"
			m := map[string]string{}
			m["FOO"] = "BAZ"
			g.Assert(s).Equal(Inject(s, m))
		})

		g.It("Should not replace vars in nil map", func() {
			s := "echo $$FOO $BAR"
			g.Assert(s).Equal(Inject(s, nil))
		})

		g.It("Should escape quoted variables", func() {
			s := `echo "$$FOO"`
			m := map[string]string{}
			m["FOO"] = "hello\nworld"
			g.Assert(`echo "hello\nworld"`).Equal(Inject(s, m))
		})

		g.It("Should replace variable prefix", func() {
			s := `tag: $${TAG=$${SHA:8}}`
			m := map[string]string{}
			m["TAG"] = ""
			m["SHA"] = "f36cbf54ee1a1eeab264c8e388f386218ab1701b"
			g.Assert("tag: f36cbf54").Equal(Inject(s, m))
		})

		g.It("Should handle nested substitution operations", func() {
			s := `echo "$${TAG##v}"`
			m := map[string]string{}
			m["TAG"] = "v1.0.0"
			g.Assert(`echo "1.0.0"`).Equal(Inject(s, m))
		})

		g.It("Should safely inject params", func() {
			m := map[string]string{
				"TOKEN":  "FOO",
				"SECRET": "BAR",
			}
			s, err := InjectSafe(before, m)
			g.Assert(err == nil).IsTrue()

			after := struct {
				Build struct {
					Image    string
					Commands []string
				}
				Deploy struct {
					Heroku struct {
						Token  string
						Secret string
					}
				}
				Publish struct {
					Amazon struct {
						Token  string
						Secret string
					}
				}
				Notify struct {
					Slack struct {
						Token  string
						Secret string
					}
				}
			}{}

			err = yaml.Unmarshal([]byte(s), &after)
			g.Assert(err == nil).IsTrue()
			g.Assert(after.Build.Commands[0]).Equal("echo $$TOKEN")
			g.Assert(after.Build.Commands[1]).Equal("echo $$SECRET")
			g.Assert(after.Deploy.Heroku.Token).Equal("FOO")
			g.Assert(after.Deploy.Heroku.Secret).Equal("BAR")
			g.Assert(after.Publish.Amazon.Token).Equal("FOO")
			g.Assert(after.Publish.Amazon.Secret).Equal("BAR")
			g.Assert(after.Notify.Slack.Token).Equal("FOO")
			g.Assert(after.Notify.Slack.Secret).Equal("BAR")
		})
	})
}

var before = `
build:
  image: foo
  commands:
    - echo $$TOKEN
    - echo $$SECRET
deploy:
  heroku:
    token: $$TOKEN
    secret: $$SECRET
publish:
  amazon:
    token: $$TOKEN
    secret: $$SECRET
notify:
  slack:
    token: $$TOKEN
    secret: $$SECRET
`
