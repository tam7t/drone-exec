package build

import (
	"encoding/json"
	"testing"

	"github.com/drone/drone-plugin-go/plugin"
)

func Test_Build(t *testing.T) {

	b, err := Parse(rawyaml)
	if err != nil {
		t.Fatal(err)
	}

	o, err := json.MarshalIndent(b.tree, " ", " ")
	if err != nil {
		t.Error(err)
	}
	println(string(o))

	s := &State{
		Repo:  plugin.Repo{FullName: "octocat/hello-world"},
		Build: plugin.Build{Commit: &plugin.Commit{Branch: "master"}},
		Job:   plugin.Job{},
	}
	b.walk(b.tree.Root, s)

}

var rawyaml = `

build:
  image: golang
  commands:
    - go build
    - go test


compose:
  redis:
    image: library/redis
  mysql: {}

publish:
  docker_x:
    name: docker1
    when:
      branch: develop
  docker_y:
    name: docker2
    when:
      branch: master
`
