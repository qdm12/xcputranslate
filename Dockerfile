ARG ALPINE_VERSION=3.12
ARG GO_VERSION=1.15

FROM --platform=linux/amd64 golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk --update add git
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/

FROM --platform=linux/amd64 base AS test
ENV CGO_ENABLED=1
RUN apk --update add g++
RUN go test -race ./...

FROM --platform=linux/amd64 base AS lint
ARG GOLANGCI_LINT_VERSION=v1.34.1
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b /usr/local/bin ${GOLANGCI_LINT_VERSION}
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM --platform=linux/amd64 base AS build
ARG TARGETPLATFORM
# The built xcputranslate binary is used
# to cross compile the program within Docker
COPY --from=qmcgaw/xcputranslate /xcputranslate /usr/local/bin/xcputranslate
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN GOARCH="$(echo ${TARGETPLATFORM} | xcputranslate -field arch)" \
  GOARM="$(echo ${TARGETPLATFORM} | xcputranslate -field arm)" \
  go build -trimpath -ldflags="-s -w \
  -X 'main.version=$VERSION' \
  -X 'main.buildDate=$BUILD_DATE' \
  -X 'main.commit=$COMMIT' \
  " -o entrypoint cmd/xcputranslate/main.go

FROM scratch
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
LABEL \
  org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
  org.opencontainers.image.created=$BUILD_DATE \
  org.opencontainers.image.version=$VERSION \
  org.opencontainers.image.revision=$COMMIT \
  org.opencontainers.image.url="https://github.com/qdm12/xcputranslate" \
  org.opencontainers.image.documentation="https://github.com/qdm12/xcputranslate" \
  org.opencontainers.image.source="https://github.com/qdm12/xcputranslate" \
  org.opencontainers.image.title="xcputranslate" \
  org.opencontainers.image.description=""
ENTRYPOINT ["/xcputranslate"]
COPY --from=build /tmp/gobuild/entrypoint /xcputranslate
