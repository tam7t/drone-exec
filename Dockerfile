# Docker image for the Drone build runner
#
#     go build
#     docker build --rm=true -t drone/drone-exec .

FROM gliderlabs/alpine:3.2
RUN apk-install ca-certificates && rm -rf /var/cache/apk/*
ADD drone-exec /bin/
ENTRYPOINT ["/bin/drone-exec"]
