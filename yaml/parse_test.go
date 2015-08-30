package yaml

import (
	"testing"

	"github.com/franela/goblin"
)

func TestParse(t *testing.T) {

	conf, err := ParseString(sample)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	g := goblin.Goblin(t)
	g.Describe("Parse Yaml", func() {

		g.It("Should parse images", func() {
			g.Assert(conf.Clone.Image).Equal("git")
		})

		g.It("Should parse image force-pull", func() {
			g.Assert(conf.Clone.Pull).Equal(true)
		})

		g.It("Should parse variable arguments", func() {
			g.Assert(conf.Clone.Vargs["path"]).Equal("github.com/octocat/hello-world")
		})

		g.It("Should parse build image", func() {
			g.Assert(conf.Build.Image).Equal("golang")
		})

		g.It("Should parse build commands", func() {
			g.Assert(conf.Build.Commands).Equal([]string{"go build", "go test"})
		})

		g.It("Should parse volume configuration", func() {
			g.Assert(conf.Build.Volumes).Equal([]string{"/tmp/volumes"})
		})

		g.It("Should parse network configuration", func() {
			g.Assert(conf.Build.Net).Equal("bridge")
		})

		g.It("Should parse environment variable map", func() {
			g.Assert(conf.Clone.Environment.Slice()).Equal(
				[]string{"GIT_DIR=.git"},
			)
		})

		g.It("Should parse environment variable slice", func() {
			g.Assert(conf.Build.Environment.Slice()).Equal(
				[]string{"GO15VENDOREXPERIMENT=1"},
			)
		})

		g.It("Should parse docker command slice", func() {
			g.Assert(conf.Compose.Slice()[0].Command.Slice()).Equal(
				[]string{
					"redis-server",
					"/usr/local/etc/redis/redis.conf",
					"--appendonly",
					"yes",
				},
			)
		})

		g.It("Should parse docker command string", func() {
			g.Assert(conf.Compose.Slice()[1].Command.Slice()).Equal(
				[]string{
					"--storageEngine",
					"wiredTiger",
				},
			)
		})

		g.It("Should allow multiple plugins of same type", func() {
			s := conf.Deploy.Slice()
			g.Assert(s[0].Image).Equal("heroku")
			g.Assert(s[1].Image).Equal("heroku")
		})

		g.It("Should maintain plugin ordering", func() {
			s := conf.Deploy.Slice()
			g.Assert(s[0].Vargs["app"]).Equal("foo.com")
			g.Assert(s[1].Vargs["app"]).Equal("dev.foo.com")
		})

		g.It("Should parse plugin filters", func() {
			s := conf.Deploy.Slice()
			g.Assert(s[0].Filter.Branch).Equal("master")
			g.Assert(s[1].Filter.Repo).Equal("octocat/helloworld")
			g.Assert(s[1].Filter.Matrix).Equal(map[string]string{"go_version": "1.5"})
		})
	})
}

var sample = `
clone:
  image: git
  pull: true
  path: github.com/octocat/hello-world
  environment:
    GIT_DIR: .git

build:
  image: golang
  environment:
    - GO15VENDOREXPERIMENT=1
  commands:
    - go build
    - go test
  volumes:
    - /tmp/volumes
  net: bridge
  privileged: true

compose:
  redis:
    image: library/redis
    command: redis-server /usr/local/etc/redis/redis.conf --appendonly yes

  mongo:
    image: library/mongo
    command:
      - --storageEngine
      - wiredTiger

deploy:
  heroku:
    app: foo.com
    when:
      branch: master
  heroku:
    app: dev.foo.com
    when:
      repo: octocat/helloworld
      branch: somebranch
      matrix:
        go_version: 1.5
`
