# XCPU Translate

A little Go static binary tool to convert Docker's buildx CPU architectures such as `linux/arm/v7` to strings for other compilers.

## Setup

```sh
VERSION=v0.1.0
ARCH=amd64
wget -O xcputranslate "https://github.com/qdm12/xcputranslate/releases/download/$VERSION/xcputranslate_$VERSION_linux_$ARCH"
chmod 500 xcputranslate
```

## Usage

```sh
echo linux/arm/v7 | xcputranslate -language golang -field arch
# arm
echo linux/arm/v7 | xcputranslate -language golang -field arm
# 7
```

More information with

```sh
xcputranslate -help
```

## Platforms supported

- `linux/amd64`
- `linux/386`
- `linux/arm64`
- `linux/arm/v6`
- `linux/arm/v7`
- `linux/s390x`
- `linux/ppc64le`

## Golang

- Use the flag `-field arch` to obtain the value to use for `GOARCH`
- Use the flag `-field arm` to obtain the value to use for `GOARM`

## Other languages

▶️ [Create an issue](https://github.com/qdm12/xcputranslate/issues/new)!

## Build it yourself

Install Go, then either

- Download it on your machine:

  ```sh
  go get github.com/qdm12/xcputranslate/cmd/xcputranslate
  ```

- Clone this repository and build it:

  ```sh
  GOARCH=arm GOARM=7 go build cmd/xcputranslate/main.go
  ```
