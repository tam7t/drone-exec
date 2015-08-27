package rules

import (
	"github.com/drone/drone-exec/yaml"
	"github.com/drone/drone-exec/yaml/rules/script"
	"github.com/drone/drone-plugin-go/plugin"
)

type RuleFunc func(*yaml.Config) error

func RuleSet(repo *plugin.Repo, build *plugin.Build, sys *plugin.System) []RuleFunc {

	var rules []RuleFunc
	rules = append(rules, PrepareClone)
	rules = append(rules, PrepareCacheRule(repo.FullName))
	rules = append(rules, PrepareImages)

	// fail the build if any sections of the yaml are
	// malformed or missing required fields.
	rules = append(rules, VerifyCache)
	rules = append(rules, VerifyBuild)
	rules = append(rules, VerifyPluginsRule(sys.Plugins[0]))
	rules = append(rules, VerifyImages)

	// remove any settings from the yaml that might
	// expose the host machine, such as volumes, network
	// and the privileged setting.
	if !isTrusted(repo, build) {
		rules = append(rules, CleanVolumes)
		rules = append(rules, CleanNetwork)
		rules = append(rules, CleanPrivileged)
	}
	// special rule that creates the build script and
	// appends to the build container.
	var wrules []script.WriterFunc
	wrules = append(wrules, script.WriteSetup())
	// wrules = append(wrules, WriteNetrc()) // machine, login, password string
	// wrules = append(wrules, WriteKey())   // key []byte
	wrules = append(wrules)
	// wrules = append(wrules, WriteCmds(nil)) // TODO requires yaml
	wrules = append(wrules, script.WriteTeardown())
	rules = append(rules, script.Rule(wrules))

	// for _, rule := range rules {
	// 	err := rule(c)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return rules
}

// helper function that returns true if a build is
// trusted, and can run in privileged mode.
func isTrusted(repo *plugin.Repo, build *plugin.Build) bool {
	return repo.Trusted && (repo.Private || isPullRequest(build))
}

// helper function that returns true if a build is
// a pull request.
func isPullRequest(build *plugin.Build) bool {
	return build.PullRequest != nil && build.PullRequest.Number != 0
}
