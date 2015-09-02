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

```sh
./drone-exec --debug --setup --build <<EOF
{
	"system": {},
	"workspace": {},
	"repo": {
		"owner": "garyburd",
		"name": "redigo",
		"full_name": "garyburd/redigo",
		"link_url": "https://github.com/garyburd/redigo",
		"clone_url": "git://github.com/garyburd/redigo.git"
	},
	"build": {
		"head_commit": {
			"ref": "refs/heads/master",
			"sha": "d8dbe4d94f15fe89232e0402c6e8a0ddf21af3ab",
			"branch": "master"
		}
	},
	"job": {
		"environment": {}
	},
	"yaml": "{ build: { image: 'golang:1.4.2', environment: ['GOPATH=/drone'], commands: ['cd redis', 'go build', 'go test -v']}, compose: { redis: { image: 'redis:2.8' } } }"
}
EOF
```

Note that the above program expects access to a Docker daemon. It will provision all the necessary build containers, execute your build, and then cleanup and remove the build environment.

### Docker

Use the following commands to build the Docker image:

```sh
# compile the binary for the correct architecture
env GOOS=linux GOARCH=amd64 go build

# build the docker image, adding the above binary
docker build --rm=true -t drone/drone-exec .
```

### Vendoring

Using the `vexp` utility to vendor dependencies:

```sh
go get https://github.com/kr/vexp
./vexp
```

