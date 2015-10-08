package runner

import (
	"testing"

	"github.com/franela/goblin"
)

func TestParse(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Yaml conditions", func() {

		g.It("Should match a branch", func() {
			g.Assert(matchBranch("master", "master")).Equal(true)
		})

		g.It("Should match a branch wildcard", func() {
			g.Assert(matchBranch("*", "master")).Equal(true)
		})

		g.It("Should match a branch with negation", func() {
			g.Assert(matchBranch("!dev", "master")).Equal(true)
		})

		g.It("Should match when branch is empty", func() {
			g.Assert(matchBranch("", "master")).Equal(true)
		})

		g.It("Should not match a branch", func() {
			g.Assert(matchBranch("dev", "master")).Equal(false)
		})

		g.It("Should not match a branch with negation", func() {
			g.Assert(matchBranch("!master", "master")).Equal(false)
		})

		g.It("Should notify on change", func() {
			g.Assert(matchChange("true", "success", "failure")).Equal(true)
			g.Assert(matchChange("true", "running", "failure")).Equal(true)
		})

		g.It("Should not notify on change when no change", func() {
			g.Assert(matchChange("true", "success", "success")).Equal(false)
		})

		g.It("Should notify on success", func() {
			g.Assert(matchSuccess("true", "success")).Equal(true)
			g.Assert(matchSuccess("true", "running")).Equal(true)
			g.Assert(matchSuccess("", "success")).Equal(true)
			g.Assert(matchSuccess("", "running")).Equal(true)
			g.Assert(matchSuccess("false", "success")).Equal(false)
			g.Assert(matchSuccess("false", "running")).Equal(false)
		})

	})

}
