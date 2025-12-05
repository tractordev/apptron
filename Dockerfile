ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG LINUX_386=linux/386
ARG LINUX_AMD64=linux/amd64
ARG GO_VERSION=1.25.0
ARG TINYGO_VERSION=0.36.0
ARG ALPINE_VERSION=3.22

FROM --platform=$LINUX_AMD64 ghcr.io/tractordev/apptron:kernel AS kernel
FROM --platform=$LINUX_AMD64 ghcr.io/progrium/v86:latest AS v86


FROM golang:$GO_VERSION-alpine AS aptn-go
WORKDIR /build
COPY system/cmd/aptn/go.mod system/cmd/aptn/go.sum ./
RUN go mod download
COPY system/cmd/aptn ./
RUN GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o aptn *.go


FROM tinygo/tinygo:$TINYGO_VERSION AS aptn-tinygo
WORKDIR /build
COPY system/cmd/aptn ./
RUN GOOS=linux GOARCH=386 tinygo build -o aptn *.go


FROM --platform=$LINUX_386 docker.io/i386/alpine:$ALPINE_VERSION AS rootfs
RUN apk add --no-cache fuse make git esbuild
COPY --from=aptn-go /build/aptn /bin/aptn
COPY ./system/bin/* /bin/
COPY ./system/etc/* /etc/


FROM alpine:$ALPINE_VERSION AS bundle-base
RUN mkdir -p /bundles
RUN apk add --no-cache brotli

FROM bundle-base AS bundle-go
ARG GO_VERSION
ENV GO_VERSION=$GO_VERSION
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-386.tar.gz \
    && tar -xzf go${GO_VERSION}.linux-386.tar.gz go/src go/pkg go/bin go/lib go/misc \
    && rm go${GO_VERSION}.linux-386.tar.gz

FROM bundle-go AS bundle-goroot
RUN tar -C /go -cf /bundles/goroot.tar . && brotli -j /bundles/goroot.tar

FROM bundle-go AS bundle-gocache-386
ENV GOCACHE=/gocache
ENV GOARCH=386
RUN /go/bin/go telemetry off && /go/bin/go build std
RUN tar -C /gocache -cf /bundles/gocache-386.tar . && brotli -j /bundles/gocache-386.tar

FROM bundle-go AS bundle-gocache-wasm
ENV GOCACHE=/gocache
ENV GOARCH=wasm
ENV GOOS=js
RUN /go/bin/go telemetry off && /go/bin/go build std
RUN tar -C /gocache -cf /bundles/gocache-wasm.tar . && brotli -j /bundles/gocache-wasm.tar

FROM bundle-base AS bundle-sys
COPY --from=rootfs / /bundle/rootfs
COPY --from=kernel /bzImage /bundle/kernel/bzImage
COPY --from=v86 /v86.wasm /bundle/v86/v86.wasm
COPY --from=v86 /bios/seabios.bin /bundle/v86/seabios.bin
COPY --from=v86 /bios/vgabios.bin /bundle/v86/vgabios.bin
RUN tar -C /bundle -czf /bundles/sys.tar.gz .


FROM golang:$GO_VERSION-alpine AS worker-build
RUN apk add --no-cache git
COPY worker/go.mod worker/go.sum ./
RUN go mod download
COPY worker .
RUN CGO_ENABLED=0 go build -o /worker ./cmd/worker


FROM scratch AS worker
COPY --from=bundle-sys /bundles/* /bundles/
COPY --from=bundle-goroot /bundles/* /bundles/
COPY --from=bundle-gocache-386 /bundles/* /bundles/
COPY --from=bundle-gocache-wasm /bundles/* /bundles/
COPY --from=worker-build /worker /worker
EXPOSE 8080
CMD ["/worker"]