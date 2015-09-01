Drone build agent that executes builds in Docker containers.

### Building

Use the following commands to build:

```sh
export GO15VENDOREXPERIMENT=1

go build
go test ./...
```

### Running

You can run the program locally for testing purposes. The build details are provided to the program via a JSON payload as seen below:

```
./drone-exec --clone --build <<EOF
{
	"system": {},
	"workspace": {},
	"repo": {
		"owner": "drone",
		"name": "drone",
		"full_name": "drone/drone",
		"link_url": "https://github.com/drone/drone",
		"clone_url": "https://github.com/drone/drone.git"
	},
	"build": {
		"head_commit": {
			"branch": "master"
		}
	},
	"job": {
		"environment": {}
	},
	"yaml": "{build: { image: golang, commands: [ go build, go test ] }, deploy: { heroku: { app: foo} }}"
}
EOF
```

Note that the above program expects access to a Docker daemon. It will provision all the necessary build containers, execute your build, and then cleanup and remove the build environment.

### Docker

Drone executes this program as a Docker container. Use the following command to build the Docker image for local integration testing within Drone:

```
# compile the binary for the correct architecture
env GOOS=linux GOARCH=amd64 go build

# build the docker image, adding the above binary
docker build --rm=true -t drone/drone-exec .
```

### Vendoring

Using the `vexp` utility to vendor dependencies:

```
go get https://github.com/kr/vexp
./vexp
```

