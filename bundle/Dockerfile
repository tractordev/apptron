ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG LINUX_386=linux/386

FROM ghcr.io/tractordev/apptron:kernel AS kernel
FROM ghcr.io/progrium/v86:latest AS v86

FROM --platform=$LINUX_386 i386/alpine:latest AS rootfs
RUN apk add -u e2fsprogs podman netavark
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