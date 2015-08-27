package main

import (
	"flag"

	"github.com/drone/drone-exec/yaml"
	"github.com/drone/drone-exec/yaml/inject"
	"github.com/drone/drone-plugin-go/plugin"

	log "github.com/Sirupsen/logrus"
)

var (
	clone  bool // execute clone steps
	build  bool // execute build steps
	deploy bool // execute deploy steps
	notify bool // execute notify steps
	debug  bool // execute in debug mode
)

func main() {

	// parses command line flags
	flag.BoolVar(&clone, "clone", false, "")
	flag.BoolVar(&build, "build", false, "")
	flag.BoolVar(&deploy, "deploy", false, "")
	flag.BoolVar(&notify, "notify", false, "")
	flag.BoolVar(&debug, "debug", false, "")
	flag.Parse()

	// parses payload options provided via
	// stdin and encoded as JSON data
	var (
		repo  = new(plugin.Repo)
		build = new(plugin.Build)
		job   = new(plugin.Job)
		sys   = new(plugin.System)
		yml   string
	)

	// TODO parse Keys
	// TODO parse Netrc

	plugin.Param("repo", repo)
	plugin.Param("build", build)
	plugin.Param("job", job)
	plugin.Param("sys", sys)
	plugin.Param("yaml", &yml)
	plugin.ParseMust()

	// configure the default log format and
	// log levels
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// injects the matrix configuration parameters
	// into the yaml prior to parsing.
	yml = inject.Inject(yml, job.Environment)

	conf, err := yaml.Parse(yml)
	if err != nil {
		log.Debugln(err) // print error messages in debug mode only
		log.Fatalln("Error parsing the .drone.yml")
	}

	log.Println(conf)

	// parse the yaml file

	// run the build
	// handle timeout, ctrl+c, etc
}
