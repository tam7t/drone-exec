package inject

import (
	"strings"
	"testing"

	"github.com/franela/goblin"
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
	})
}

var before = `
build:
  image: foo
  commands:
    - echo $$TOKEN
    - echo $$SECRET
deploy:
  digital_ocean:
    token: $$TOKEN
    secret: $$SECRET
`

var after = `
cache: {}
clone: {}
build:
  image: foo
  commands:
  - echo $$TOKEN
  - echo $$SECRET
compose: {}
publish: {}
deploy:
  digital_ocean:
    token: FOO
    secret: BAR
notify: {}
`
