package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	builder "github.com/drone/drone-exec/build"
	"github.com/drone/drone-exec/build/parse"
	"github.com/drone/drone-exec/docker"
	"github.com/drone/drone-exec/yaml/inject"
	"github.com/drone/drone-exec/yaml/path"
	"github.com/drone/drone-plugin-go/plugin"
	"github.com/samalba/dockerclient"

	log "github.com/Sirupsen/logrus"
)

var (
	setup  bool // execute clone steps
	build  bool // execute build steps
	deploy bool // execute deploy steps
	notify bool // execute notify steps
	debug  bool // execute in debug mode
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
	flag.BoolVar(&setup, "setup", false, "")
	flag.BoolVar(&build, "build", false, "")
	flag.BoolVar(&deploy, "deploy", false, "")
	flag.BoolVar(&notify, "notify", false, "")
	flag.BoolVar(&debug, "debug", false, "")
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

	// extracts the clone path from the yaml. If
	// the clone path doesn't exist it uses a path
	// derrived from the repository uri.
	payload.Workspace.Path = path.Parse(yml, payload.Repo.Link)
	payload.Workspace.Root = "/drone/src"

	b, err := builder.Parse(yml)
	if err != nil {
		log.Debugln(err) // print error messages in debug mode only
		log.Fatalln("Error parsing the .drone.yml")
		os.Exit(1)
	}

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

	state := &builder.State{
		Client:    controller,
		Stdout:    os.Stdout,
		Stderr:    os.Stdout,
		Repo:      payload.Repo,
		Build:     payload.Build,
		Job:       payload.Job,
		System:    payload.System,
		Workspace: payload.Workspace,
	}
	if setup {
		err = b.RunNode(state, parse.NodeCache|parse.NodeClone)
		if err != nil {
			log.Debugln(err)
		}
	}
	if build && !state.IsFailed() {
		err = b.RunNode(state, parse.NodeCompose|parse.NodeBuild)
		if err != nil {
			log.Debugln(err)
		}
	}
	if deploy && !state.IsFailed() {
		err = b.RunNode(state, parse.NodePublish|parse.NodeDeploy)
		if err != nil {
			log.Debugln(err)
		}
	}
	if setup {
		err = b.RunNode(state, parse.NodeCache)
		if err != nil {
			log.Debugln(err)
		}
	}
	if notify {
		err = b.RunNode(state, parse.NodeNotify)
		if err != nil {
			log.Debugln(err)
		}
	}

	if state.IsFailed() {
		controller.Destroy()
		os.Exit(state.ExitCode())
	}
}

// applies verification and transformation rules
// err = apply.Rules(&apply.Context{conf, repo, build, job, sys, ws})
// if err != nil {
// 	log.Debugln(err) // print error messages in debug mode only
// 	log.Fatalln("Error linting the .drone.yml")
// }

// // and append the matrix parameters as environment
// // variables for the build
// for k, v := range c.Job.Environment {
// 	env := k + "=" + v
// 	c.Conf.Build.Environment = append(c.Conf.Build.Environment, env)
// }

// // and append drone, jenkins, travis and other
// // environment variables that may be of use.
// for k, v := range toEnv(c) {
// 	env := k + "=" + v
// 	c.Conf.Build.Environment = append(c.Conf.Build.Environment, env)
// }
