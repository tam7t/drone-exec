package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/drone/drone-exec/docker"
	"github.com/drone/drone-exec/parser"
	"github.com/drone/drone-exec/runner"
	"github.com/drone/drone-exec/yaml/inject"
	"github.com/drone/drone-exec/yaml/path"
	"github.com/drone/drone-plugin-go/plugin"
	"github.com/samalba/dockerclient"

	log "github.com/Sirupsen/logrus"
)

var (
	cache  bool // execute cache steps
	clone  bool // execute clone steps
	build  bool // execute build steps
	deploy bool // execute deploy steps
	notify bool // execute notify steps
	debug  bool // execute in debug mode
	force  bool // force pull plugin images
)

// payload defines the raw plugin payload that
// stores the build metadata and configuration.
var payload = struct {
	Yaml      string            `json:"yaml"`
	Repo      *plugin.Repo      `json:"repo"`
	Build     *plugin.Build     `json:"build"`
	Job       *plugin.Job       `json:"job"`
	System    *plugin.System    `json:"system"`
	Workspace *plugin.Workspace `json:"workspace"`
}{}

func main() {

	// parses command line flags
	flag.BoolVar(&cache, "cache", false, "")
	flag.BoolVar(&clone, "clone", false, "")
	flag.BoolVar(&build, "build", false, "")
	flag.BoolVar(&deploy, "deploy", false, "")
	flag.BoolVar(&notify, "notify", false, "")
	flag.BoolVar(&debug, "debug", false, "")
	flag.BoolVar(&force, "pull", false, "")
	flag.Parse()

	// unmarshal the json payload via stdin or
	// via the command line args (whichever was used)
	plugin.MustUnmarshal(&payload)

	// configure the default log format and
	// log levels
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// injects the matrix configuration parameters
	// into the yaml prior to parsing.
	yml := inject.Inject(payload.Yaml, payload.Job.Environment)
	yml = inject.Inject(yml, map[string]string{
		"COMMIT":       payload.Build.Commit.Sha,
		"BRANCH":       payload.Build.Commit.Branch,
		"BUILD_NUMBER": strconv.Itoa(payload.Build.Number),
	})

	// extracts the clone path from the yaml. If
	// the clone path doesn't exist it uses a path
	// derrived from the repository uri.
	payload.Workspace.Path = path.Parse(yml, payload.Repo.Link)
	payload.Workspace.Root = "/drone/src"

	rules := []parser.RuleFunc{
		parser.ImageName,
		parser.ImageMatchFunc(payload.System.Plugins),
		parser.ImagePullFunc(force),
		parser.SanitizeFunc(payload.Repo.Trusted),
		parser.CacheFunc(payload.Repo.FullName),
		parser.Escalate,
	}
	tree, err := parser.Parse(yml, rules)
	if err != nil {
		log.Debugln(err) // print error messages in debug mode only
		log.Fatalln("Error parsing the .drone.yml")
		os.Exit(1)
	}
	r := runner.Load(tree)

	client, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		log.Debugln(err)
		log.Fatalln("Error creating the docker client.")
		os.Exit(1)
	}

	// // creates a wrapper Docker client that uses an ambassador
	// // container to create a pod-like environment.
	controller, err := docker.NewClient(client)
	if err != nil {
		log.Debugln(err)
		log.Fatalln("Error creating the docker ambassador.")
		os.Exit(1)
	}
	defer controller.Destroy()

	// watch for sigkill (timeout or cancel build)
	killc := make(chan os.Signal, 1)
	signal.Notify(killc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-killc
		log.Println("Cancel request received, killing process")
		controller.Destroy() // possibe race here. implement lock on the other end
		os.Exit(130)         // cancel is treated like ctrl+c
	}()

	go func() {
		var timeout = payload.Repo.Timeout
		if timeout == 0 {
			timeout = 60
		}
		<-time.After(time.Duration(timeout) * time.Minute)
		log.Println("Timeout request received, killing process")
		controller.Destroy() // possibe race here. implement lock on the other end
		os.Exit(128)         // cancel is treated like ctrl+c
	}()

	state := &runner.State{
		Client:    controller,
		Stdout:    os.Stdout,
		Stderr:    os.Stdout,
		Repo:      payload.Repo,
		Build:     payload.Build,
		Job:       payload.Job,
		System:    payload.System,
		Workspace: payload.Workspace,
	}
	if cache {
		err = r.RunNode(state, parser.NodeCache)
		if err != nil {
			log.Debugln(err)
		}
	}
	if clone {
		err = r.RunNode(state, parser.NodeClone)
		if err != nil {
			log.Debugln(err)
		}
	}
	if build && !state.Failed() {
		err = r.RunNode(state, parser.NodeCompose|parser.NodeBuild)
		if err != nil {
			log.Debugln(err)
		}
	}
	if deploy && !state.Failed() {
		err = r.RunNode(state, parser.NodePublish|parser.NodeDeploy)
		if err != nil {
			log.Debugln(err)
		}
	}
	if cache {
		err = r.RunNode(state, parser.NodeCache)
		if err != nil {
			log.Debugln(err)
		}
	}
	if notify {
		err = r.RunNode(state, parser.NodeNotify)
		if err != nil {
			log.Debugln(err)
		}
	}

	if state.Failed() {
		controller.Destroy()
		os.Exit(state.ExitCode())
	}
}
