package docker

import (
	"errors"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

var (
	ErrTimeout = errors.New("Timeout")
	ErrLogging = errors.New("Logs not available")
)

var (
	// options to fetch the stdout and stderr logs
	logOpts = &dockerclient.LogOptions{
		Stdout: true,
		Stderr: true,
	}

	// options to fetch the stdout and stderr logs
	// by tailing the output.
	logOptsTail = &dockerclient.LogOptions{
		Follow: true,
		Stdout: true,
		Stderr: true,
	}
)

func Run(client dockerclient.Client, conf *dockerclient.ContainerConfig, pull bool) (*dockerclient.ContainerInfo, error) {

	// force-pull the image if specified.
	// TEMPORARY: always try to pull the new image for now
	// since we'll be frequently updating the plugin images
	// over the next few weeks
	if pull || strings.HasPrefix(conf.Image, "plugins/") {
		client.PullImage(conf.Image, nil)
	}

	// attempts to create the contianer
	id, err := client.CreateContainer(conf, "")
	if err != nil {
		// and pull the image and re-create if that fails
		err = client.PullImage(conf.Image, nil)
		if err != nil {
			log.Errorf("Error pulling %s. %s\n", conf.Image, err)
			return nil, err
		}

		id, err = client.CreateContainer(conf, "")
		// make sure the container is removed in
		// the event of a creation error.
		if err != nil {
			log.Errorf("Error starting %s. %s\n", conf.Image, err)
			client.RemoveContainer(id, true, true)
			return nil, err
		}
	}

	// ensures the container is always stopped
	// and ready to be removed.
	defer func() {
		client.StopContainer(id, 5)
		client.KillContainer(id, "9")
	}()

	// fetches the container information.
	info, err := client.InspectContainer(id)
	if err != nil {
		log.Errorf("Error inspecting %s. %s\n", conf.Image, err)
		client.RemoveContainer(id, true, true)
		return nil, err
	}

	// channel listening for errors while the
	// container is running async.
	errc := make(chan error, 1)
	infoc := make(chan *dockerclient.ContainerInfo, 1)
	go func() {

		// starts the container
		err := client.StartContainer(id, &conf.HostConfig)
		if err != nil {
			log.Errorf("Error starting %s. %s\n", conf.Image, err)
			errc <- err
			return
		}

		// blocks and waits for the container to finish
		// by streaming the logs (to /dev/null). Ideally
		// we could use the `wait` function instead
		rc, err := client.ContainerLogs(id, logOptsTail)
		if err != nil {
			log.Errorf("Error tailing %s. %s\n", conf.Image, err)
			errc <- err
			return
		}
		defer rc.Close()
		StdCopy(os.Stdout, os.Stdout, rc)

		// fetches the container information
		info, err := client.InspectContainer(id)
		if err != nil {
			log.Errorf("Error getting exit code for %s. %s\n", conf.Image, err)
			errc <- err
			return
		}
		infoc <- info
	}()

	select {
	case info := <-infoc:
		return info, nil
	case err := <-errc:
		return info, err
		// TODO checkout net.Context and cancel
		// case <-time.After(timeout):
		// 	return info, ErrTimeout
	}
}

func RunDaemon(client dockerclient.Client, conf *dockerclient.ContainerConfig, pull bool) (*dockerclient.ContainerInfo, error) {
	// force-pull the image if specified.
	// TEMPORARY: always try to pull the new image for now
	// since we'll be frequently updating the plugin images
	// over the next few weeks
	if pull || strings.HasPrefix(conf.Image, "plugins/") {
		client.PullImage(conf.Image, nil)
	}

	// attempts to create the contianer
	id, err := client.CreateContainer(conf, "")
	if err != nil {
		// and pull the image and re-create if that fails
		err = client.PullImage(conf.Image, nil)
		if err != nil {
			log.Errorf("Error pulling %s. %s\n", conf.Image, err)
			return nil, err
		}
		id, err = client.CreateContainer(conf, "")
		if err != nil {
			log.Errorf("Error creating %s. %s\n", conf.Image, err)
			client.RemoveContainer(id, true, true)
			return nil, err
		}
	}

	// fetches the container information
	info, err := client.InspectContainer(id)
	if err != nil {
		log.Errorf("Error inspecting %s. %s\n", conf.Image, err)
		client.RemoveContainer(id, true, true)
		return nil, err
	}

	// starts the container
	err = client.StartContainer(id, &conf.HostConfig)
	if err != nil {
		log.Errorf("Error starting %s. %s\n", conf.Image, err)
	}
	return info, err
}
