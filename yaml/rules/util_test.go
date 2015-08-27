package rules

import (
	"testing"

	"github.com/franela/goblin"
)

func Test_Utils(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Rules", func() {

		g.It("Should return full qualified plugin name", func() {
			g.Assert(pluginName("microsoft/azure")).Equal("microsoft/azure")
			g.Assert(pluginName("azure")).Equal("plugins/drone-azure")
			g.Assert(pluginName("azure_storage")).Equal("plugins/drone-azure-storage")
		})

		g.It("Should return full qualified plugin name if no default", func() {
			g.Assert(pluginNameDefault("hg", "git")).Equal("plugins/drone-hg")
			g.Assert(pluginNameDefault("", "git")).Equal("plugins/drone-git")
		})
	})
}
