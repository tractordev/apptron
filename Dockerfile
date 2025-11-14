ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG LINUX_386=linux/386
ARG LINUX_AMD64=linux/amd64
ARG GO_VERSION=1.25.0

FROM --platform=$LINUX_AMD64 ghcr.io/tractordev/apptron:kernel AS kernel
FROM --platform=$LINUX_AMD64 ghcr.io/progrium/v86:latest AS v86

FROM golang:$GO_VERSION-alpine AS wexec-go
WORKDIR /build
COPY system/cmd/wexec/go.mod system/cmd/wexec/go.sum ./
RUN go mod download
COPY system/cmd/wexec ./
RUN GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o wexec *.go

FROM tinygo/tinygo:0.36.0 AS wexec-tinygo
WORKDIR /build
COPY system/cmd/wexec ./
RUN GOOS=linux GOARCH=386 tinygo build -o wexec *.go

FROM --platform=$LINUX_386 docker.io/i386/alpine:latest AS rootfs
RUN apk add --no-cache fuse
COPY --from=wexec-go /build/wexec /bin/wexec
COPY ./system/bin/* /bin/
COPY ./system/etc/* /etc/

FROM alpine:3.22 AS bundle
ARG GO_VERSION
ENV GO_VERSION=$GO_VERSION
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-386.tar.gz \
    && tar -xzf go${GO_VERSION}.linux-386.tar.gz go/src go/pkg go/bin go/lib go/misc \
    && rm go${GO_VERSION}.linux-386.tar.gz
ENV GOCACHE=/go/cache
RUN /go/bin/go telemetry off && /go/bin/go build std
RUN tar -C /go -czf /gobundle.tgz .
COPY --from=rootfs / /bundle/rootfs
COPY --from=kernel /bzImage /bundle/kernel/bzImage
COPY --from=v86 /v86.wasm /bundle/v86/v86.wasm
COPY --from=v86 /bios/seabios.bin /bundle/v86/seabios.bin
COPY --from=v86 /bios/vgabios.bin /bundle/v86/vgabios.bin
RUN tar -C /bundle -czf /bundle.tgz .

FROM golang:1.24.5-alpine AS golang-build
WORKDIR /build
RUN apk add --no-cache git

FROM golang-build AS session-build
COPY session/go.mod session/go.sum ./
RUN go mod download
COPY session/main.go ./
RUN CGO_ENABLED=0 go build -o /session .

FROM scratch AS session
COPY --from=bundle /bundle.tgz /bundle.tgz
COPY --from=bundle /gobundle.tgz /gobundle.tgz
COPY --from=session-build /session /session
EXPOSE 8080
CMD ["/session"]