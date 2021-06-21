# XCPU Translate

A little Go static binary tool to convert Docker's buildx CPU architectures such as `linux/arm/v7` to strings for other compilers.

🆕 Sleep before building to prevent build out of memory issues, depending on the target platform. See [moby/buildkit#1131](https://github.com/moby/buildkit/issues/1131) for more context.

## Setup and usage

💡 It should be used with at least one of:

- Docker [buildx](https://docs.docker.com/buildx/working-with-buildx/) builds
- [BuildKit](https://github.com/moby/buildkit) directly
- Docker builds running with [`DOCKER_BUILDKIT=1`](https://docs.docker.com/develop/develop-images/build_enhancements/#to-enable-buildkit-builds)

### Docker platform translation

The following shows an example on how to use it to cross compile a Go program.

We compile for `linux/arm/v7` on a `linux/amd64` machine using:

```sh
docker build --platform linux/arm/v7 .
```

```Dockerfile
# Note you cannot COPY directly from the image or it will duplicate instructions
# for each target platform. You need to FROM it and then COPY from the alias.
FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:v0.6.0 AS xcputranslate

# We use the builder native architecture to build the program
FROM --from=${BUILDPLATFORM} golang:1.16-alpine3.13 AS build
# The build argument TARGETPLATFORM is automatically
# plugged in by docker build
ARG TARGETPLATFORM

# Setup additional build dependencies
RUN apk --update add git
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
# Download your Go modules
COPY go.mod go.sum ./
RUN go mod download
# Copy your source code
COPY . .

# 📥 Install xcputranslate for your build architecture
COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate

# 🦾 We cross build for linux/arm/v7
RUN GOARCH="$(xcputranslate translate -targetplatform ${TARGETPLATFORM} -language golang -field arch)" \
    GOARM="$(xcputranslate translate -targetplatform ${TARGETPLATFORM} -language golang -field arm)" \
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
COPY --from=qmcgaw/xcputranslate:v0.6.0 /xcputranslate /usr/local/bin/xcputranslate
```

### Sequential cross CPU Docker builds

For now, Buildkit will run all your target platform specific build instructions (such as `go build`) in parallel. This can be nice but can also cause out of memory errors, even on CIs such as Github Actions. I had the problem with a `go build` cross compiling for 5+ architectures in parallel. See [moby/buildkit#1131](https://github.com/moby/buildkit/issues/1131) for more context.

To fix this temporay problem, `xcputranslate` extends its feature by adding a new command: `xcputranslate sleep`

It allows to sleep before building depending on the target platform and a list of target platforms.

For example:

```Dockerfile
ARG ALLTARGETPLATFORMS=linux/amd64,linux/386
RUN xcputranslate sleep -targetplatform=${TARGETPLATFORM} -order=${ALLTARGETPLATFORMS} && \
    GOARCH="$(xcputranslate translate -targetplatform ${TARGETPLATFORM} -language golang -field arch)" \
    GOARM="$(xcputranslate translate -targetplatform ${TARGETPLATFORM} -language golang -field arm)" \
    go build -o entrypoint main.go
```

will sleep 0 for `linux/amd64` and 3 seconds for `linux/386`.

The `-order` defaults to `linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6,linux/386,linux/ppc64le,linux/s390x,linux/riscv64` which is an order sorted by popularity. It means that, for example, for `-targetplatform=linux/arm/v7`, it will sleep `buildtime x 2` where 2 is the order index of the target platform.

The `-buildtime` flag allows to set the estimated build time.

If for example, your build takes 15 seconds and you want to target the platforms `linux/arm64` and `linux/s390x` only, you should set `-buildtime=15s -order=linux/arm64,linux/s390x` to have the shortest sleep times possible and each build run sequentially.

### Out of Docker

You can also run already built binaries out of Docker:

```sh
# Install
VERSION=v0.6.0
ARCH=amd64
wget -O xcputranslate "https://github.com/qdm12/xcputranslate/releases/download/$VERSION/xcputranslate_$VERSION_linux_$ARCH"
chmod +x xcputranslate

# Run
xcputranslate translate -targetplatform "linux/arm/v7" -language golang -field arch
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
