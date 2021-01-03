# XCPU Translate

A little Go static binary tool to convert Docker's buildx CPU architectures such as `linux/arm/v7` to strings for other compilers.

## Usage

```sh
echo linux/arm/v7 | xcputranslate -language golang -field arm
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
