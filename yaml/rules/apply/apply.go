package apply

import (
	"github.com/drone/drone-exec/yaml"
	"github.com/drone/drone-exec/yaml/rules"
	"github.com/drone/drone-exec/yaml/rules/script"
	"github.com/drone/drone-plugin-go/plugin"
)

type Context struct {
	Conf  *yaml.Config
	Repo  *plugin.Repo
	Netrc *plugin.Netrc
	Keys  *plugin.Keypair
	Build *plugin.Build
	Job   *plugin.Job
	Sys   *plugin.System
}

func Rules(c *Context) error {

	var funcs []rules.RuleFunc
	funcs = append(funcs, rules.PrepareClone)
	funcs = append(funcs, rules.PrepareCacheRule(c.Repo.FullName))
	funcs = append(funcs, rules.PrepareImages)

	// fail the build if any sections of the yaml are
	// malformed or missing required fields.
	funcs = append(funcs, rules.LintCache)
	funcs = append(funcs, rules.LintBuild)
	funcs = append(funcs, rules.LintPluginsRule(c.Sys.Plugins))
	funcs = append(funcs, rules.LintImages)

	// remove any settings from the yaml that might
	// expose the host machine, such as volumes, network
	// and the privileged setting.
	if !isTrusted(c.Repo, c.Build) {
		funcs = append(funcs, rules.CleanVolumes)
		funcs = append(funcs, rules.CleanNetwork)
		funcs = append(funcs, rules.CleanPrivileged)
	}

	// enable certain white-listed plugins to run
	// in privileged mode (ie Docker)
	funcs = append(funcs, rules.EnableDocker)

	// special rule that creates the build script and
	// appends to the build container.
	var writers []script.WriterFunc
	writers = append(writers, script.WriteSetup())

	// only include the ntrc and key in the build
	// environment if the repository is private.
	if c.Repo.Private && c.Netrc != nil && c.Keys != nil {
		writers = append(writers, script.WriteNetrc(c.Netrc.Machine, c.Netrc.Login, c.Netrc.Password))
		writers = append(writers, script.WriteKey([]byte(c.Keys.Private)))
	}

	writers = append(writers, script.WriteCmds(c.Conf))
	writers = append(writers, script.WriteTeardown())
	funcs = append(funcs, script.RuleFunc(writers))

	for _, rule := range funcs {
		err := rule(c.Conf)
		if err != nil {
			return err
		}
	}

	return nil
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
