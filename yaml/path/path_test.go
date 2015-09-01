package path

import (
	"testing"

	"github.com/franela/goblin"
)

func Test_Parse(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Parsing the yaml", func() {

		g.It("Should return the clone path", func() {
			p := Parse(sampleClone, "http://github.com/foo/bar")
			g.Assert(p).Equal("/drone/src/github.com/octocat/hello-world")
		})

		g.It("Should handle missing clone path", func() {
			p := Parse(sampleEmpty, "http://github.com/foo/bar")
			g.Assert(p).Equal("/drone/src/github.com/foo/bar")
		})

		g.It("Should handle mission clone section", func() {
			p := Parse(sampleMissing, "http://github.com/foo/bar")
			g.Assert(p).Equal("/drone/src/github.com/foo/bar")
		})

		g.It("Should exclude port numbers from the url", func() {
			p := Parse(sampleMissing, "http://github.com:80/foo/bar")
			g.Assert(p).Equal("/drone/src/github.com/foo/bar")
		})

		g.It("Should not prepend root if already part of path", func() {
			p := Parse(sampleAbs, "http://github.com/foo/bar")
			g.Assert(p).Equal("/drone/src/github.com/octocat/hello-world")
		})

		g.It("Should use an empty path when the url is malformed", func() {
			p := Parse(sampleMissing, "%gh&%ij")
			g.Assert(p).Equal("/drone/src")
		})
	})
}

var sampleClone = `
clone:
  path: github.com/octocat/hello-world
`

var sampleEmpty = `
clone: {}
`

var sampleMissing = `
build: {}
`

var sampleAbs = `
clone:
  path: /drone/src/github.com/octocat/hello-world
`
