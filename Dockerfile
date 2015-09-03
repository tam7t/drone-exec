# Docker image for the Drone build runner
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t drone/drone-exec .

FROM gliderlabs/alpine:3.1
RUN apk-install ca-certificates && rm -rf /var/cache/apk/*
ADD drone-exec /bin/
ENTRYPOINT ["/bin/drone-exec"]
