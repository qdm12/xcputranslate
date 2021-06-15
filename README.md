# XCPU Translate

A little Go static binary tool to convert Docker's buildx CPU architectures such as `linux/arm/v7` to strings for other compilers.

## Setup and usage

💡 It should be used with either:

- Docker [buildx](https://docs.docker.com/buildx/working-with-buildx/) builds
- Docker builds running with [`DOCKER_BUILDKIT=1`](https://docs.docker.com/develop/develop-images/build_enhancements/#to-enable-buildkit-builds)

The following shows an example on how to use it to cross compile a Go program.

We compile for `linux/arm/v7` on a `linux/amd64` machine using:

```sh
docker build --platform linux/arm/v7 .
```

```Dockerfile
# We use the builder native architecture to build the program
FROM --from=${BUILDPLATFORM} golang:1.16-alpine3.13 AS build
# The build argument TARGETPLATFORM is automatically
# plugged in by docker build
ARG TARGETPLATFORM

# 📥 Install xcputranslate for your build architecture
COPY --from=qmcgaw/xcputranslate /xcputranslate /usr/local/bin/xcputranslate

# Setup additional build dependencies
RUN apk --update add git
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
# Download your Go modules
COPY go.mod go.sum ./
RUN go mod download
# Copy your source code
COPY . .

# 🦾 We cross build for linux/arm/v7
RUN GOARCH="$(xcputranslate -targetplatform ${TARGETPLATFORM} -language golang -field arch)" \
    GOARM="$(xcputranslate -targetplatform ${TARGETPLATFORM} -language golang -field arm)" \
    go build -o entrypoint main.go

# This is built on the target architecture (e.g. linux/arm/v7)
FROM alpine:3.13
# Run as user ID 1000, not the default root
USER 1000
ENTRYPOINT ["/usr/local/bin/entrypoint"]
COPY --from=build --chown=1000 /tmp/gobuild/entrypoint /usr/local/bin/entrypoint
```

Note that you can also specify a Docker tag to have the program matching a certain Github release. For example:

```Dockerfile
COPY --from=qmcgaw/xcputranslate:v0.4.0 /xcputranslate /usr/local/bin/xcputranslate
```

### Out of Docker

You can also run already built binaries out of Docker:

```sh
# Install
VERSION=v0.4.0
ARCH=amd64
wget -O xcputranslate "https://github.com/qdm12/xcputranslate/releases/download/$VERSION/xcputranslate_$VERSION_linux_$ARCH"
chmod +x xcputranslate

# Run
xcputranslate -targetplatform "linux/arm/v7" -language golang -field arch
# 7
```

## Docker platforms supported

- `linux/amd64`
- `linux/386`
- `linux/arm64` & `linux/arm64/v8`
- `linux/arm/v6`
- `linux/arm/v7`
- `linux/s390x`
- `linux/ppc64le`
- `linux/riscv64`

## Languages supported

### Golang

- Use the flag `-field arch` to obtain the value to use for `GOARCH`
- Use the flag `-field arm` to obtain the value to use for `GOARM`

### Uname

Not really a language, although it gives the same as `uname -m` on Linux OSes.
For example `linux/arm64` gives `aarch64`. This is useful for Rust commands for example.

Use it using `-language=uname` and with `-field arch`.

### Other languages

▶️ [Create an issue](https://github.com/qdm12/xcputranslate/issues/new)!
