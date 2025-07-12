# syntax = docker/dockerfile-upstream:1.15.1-labs

# THIS FILE WAS AUTOMATICALLY GENERATED, PLEASE DO NOT EDIT.
#
# Generated on 2025-05-29T16:07:21Z by kres 9f64b0d.

ARG TOOLCHAIN

FROM ghcr.io/siderolabs/ca-certificates:v1.10.0 AS image-ca-certificates

FROM ghcr.io/siderolabs/fhs:v1.10.0 AS image-fhs

# collects proto specs
FROM scratch AS proto-specs
ADD https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto /api/google/api/
ADD https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto /api/google/api/
ADD api/logs/logs.proto /api/logs/

# base toolchain image
FROM --platform=${BUILDPLATFORM} ${TOOLCHAIN} AS toolchain
RUN apk --update --no-cache add bash curl build-base protoc protobuf-dev

# build tools
FROM --platform=${BUILDPLATFORM} toolchain AS tools
ENV GO111MODULE=on
ARG CGO_ENABLED
ENV CGO_ENABLED=${CGO_ENABLED}
ARG GOTOOLCHAIN
ENV GOTOOLCHAIN=${GOTOOLCHAIN}
ARG GOEXPERIMENT
ENV GOEXPERIMENT=${GOEXPERIMENT}
ENV GOPATH=/go
ARG GOIMPORTS_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go install golang.org/x/tools/cmd/goimports@v${GOIMPORTS_VERSION}
RUN mv /go/bin/goimports /bin
ARG GOMOCK_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go install go.uber.org/mock/mockgen@v${GOMOCK_VERSION}
RUN mv /go/bin/mockgen /bin
ARG PROTOBUF_GO_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOBUF_GO_VERSION}
RUN mv /go/bin/protoc-gen-go /bin
ARG GRPC_GO_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${GRPC_GO_VERSION}
RUN mv /go/bin/protoc-gen-go-grpc /bin
ARG GRPC_GATEWAY_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION}
RUN mv /go/bin/protoc-gen-grpc-gateway /bin
ARG VTPROTOBUF_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@v${VTPROTOBUF_VERSION}
RUN mv /go/bin/protoc-gen-go-vtproto /bin
ARG DEEPCOPY_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go install github.com/siderolabs/deep-copy@${DEEPCOPY_VERSION} \
	&& mv /go/bin/deep-copy /bin/deep-copy
ARG GOLANGCILINT_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCILINT_VERSION} \
	&& mv /go/bin/golangci-lint /bin/golangci-lint
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go install golang.org/x/vuln/cmd/govulncheck@latest \
	&& mv /go/bin/govulncheck /bin/govulncheck
ARG GOFUMPT_VERSION
RUN go install mvdan.cc/gofumpt@${GOFUMPT_VERSION} \
	&& mv /go/bin/gofumpt /bin/gofumpt

# tools and sources
FROM tools AS base
WORKDIR /src
COPY go.mod go.mod
COPY go.sum go.sum
RUN cd .
RUN --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go mod download
RUN --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go mod verify
COPY ./api ./api
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg go list -mod=readonly all >/dev/null

# runs protobuf compiler
FROM tools AS proto-compile
COPY --from=proto-specs / /
RUN protoc -I/api --grpc-gateway_out=paths=source_relative:/api --grpc-gateway_opt=generate_unbound_methods=true --grpc-gateway_opt=standalone=true /api/google/api/http.proto /api/google/api/annotations.proto
RUN protoc -I/api --grpc-gateway_out=paths=source_relative:/api --grpc-gateway_opt=generate_unbound_methods=true --go_out=paths=source_relative:/api --go-grpc_out=paths=source_relative:/api --go-vtproto_out=paths=source_relative:/api --go-vtproto_opt=features=marshal+unmarshal+size+equal+clone /api/logs/logs.proto
RUN rm /api/logs/logs.proto
RUN goimports -w -local github.com/unix4ever/fe-test-service /api
RUN gofumpt -w /api

# runs gofumpt
FROM base AS lint-gofumpt
RUN FILES="$(gofumpt -l .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'gofumpt -w .':\n${FILES}"; exit 1)

# runs golangci-lint
FROM base AS lint-golangci-lint
WORKDIR /src
COPY .golangci.yml .
ENV GOGC=50
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/root/.cache/golangci-lint,id=fe-test-service/root/.cache/golangci-lint,sharing=locked --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg golangci-lint run --config .golangci.yml

# runs govulncheck
FROM base AS lint-govulncheck
WORKDIR /src
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg govulncheck ./...

# runs unit-tests with race detector
FROM base AS unit-tests-race
WORKDIR /src
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg --mount=type=cache,target=/tmp,id=fe-test-service/tmp CGO_ENABLED=1 go test -v -race -count 1 ${TESTPKGS}

# runs unit-tests
FROM base AS unit-tests-run
WORKDIR /src
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg --mount=type=cache,target=/tmp,id=fe-test-service/tmp go test -v -covermode=atomic -coverprofile=coverage.txt -coverpkg=${TESTPKGS} -count 1 ${TESTPKGS}

# cleaned up specs and compiled versions
FROM scratch AS generate
COPY --from=proto-compile /api/ /api/

FROM scratch AS unit-tests
COPY --from=unit-tests-run /src/coverage.txt /coverage-unit-tests.txt

# builds server-darwin-amd64
FROM base AS server-darwin-amd64-build
COPY --from=generate / /
WORKDIR /src/cmd/server
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg GOARCH=amd64 GOOS=darwin go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS}" -o /server-darwin-amd64

# builds server-darwin-arm64
FROM base AS server-darwin-arm64-build
COPY --from=generate / /
WORKDIR /src/cmd/server
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg GOARCH=arm64 GOOS=darwin go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS}" -o /server-darwin-arm64

# builds server-linux-amd64
FROM base AS server-linux-amd64-build
COPY --from=generate / /
WORKDIR /src/cmd/server
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg GOARCH=amd64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS}" -o /server-linux-amd64

# builds server-linux-arm64
FROM base AS server-linux-arm64-build
COPY --from=generate / /
WORKDIR /src/cmd/server
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=fe-test-service/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=fe-test-service/go/pkg GOARCH=arm64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS}" -o /server-linux-arm64

FROM scratch AS server-darwin-amd64
COPY --from=server-darwin-amd64-build /server-darwin-amd64 /server-darwin-amd64

FROM scratch AS server-darwin-arm64
COPY --from=server-darwin-arm64-build /server-darwin-arm64 /server-darwin-arm64

FROM scratch AS server-linux-amd64
COPY --from=server-linux-amd64-build /server-linux-amd64 /server-linux-amd64

FROM scratch AS server-linux-arm64
COPY --from=server-linux-arm64-build /server-linux-arm64 /server-linux-arm64

FROM server-linux-${TARGETARCH} AS server

FROM scratch AS server-all
COPY --from=server-darwin-amd64 / /
COPY --from=server-darwin-arm64 / /
COPY --from=server-linux-amd64 / /
COPY --from=server-linux-arm64 / /

FROM scratch AS image-server
ARG TARGETARCH
COPY --from=server server-linux-${TARGETARCH} /server
COPY --from=image-fhs / /
COPY --from=image-ca-certificates / /
LABEL org.opencontainers.image.source=https://github.com/unix4ever/fe-test-service
ENTRYPOINT ["/server"]

