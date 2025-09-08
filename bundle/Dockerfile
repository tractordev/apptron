ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG LINUX_386=linux/386
ARG LINUX_AMD64=linux/amd64

FROM --platform=$LINUX_AMD64 ghcr.io/tractordev/apptron:kernel AS kernel
FROM --platform=$LINUX_AMD64 ghcr.io/progrium/v86:latest AS v86

FROM --platform=$TARGETPLATFORM docker.io/i386/alpine:latest AS rootfs
COPY ./bin/* /bin/
COPY ./etc/* /etc/

FROM scratch AS bundle
COPY --from=rootfs / /rootfs
COPY --from=kernel /bzImage /kernel/bzImage
COPY --from=v86 /libv86.js /v86/libv86.js
COPY --from=v86 /v86.wasm /v86/v86.wasm
COPY --from=v86 /bios/seabios.bin /v86/seabios.bin
COPY --from=v86 /bios/vgabios.bin /v86/vgabios.bin
COPY ./init.js /init.js
CMD ["true"]