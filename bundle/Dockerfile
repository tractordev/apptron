ARG BUILDPLATFORM
ARG TARGETPLATFORM

FROM --platform=$TARGETPLATFORM i386/alpine:latest AS rootfs
RUN apk add -u e2fsprogs podman netavark
COPY ./bin/* /bin/
COPY ./etc/* /etc/

FROM --platform=$BUILDPLATFORM ghcr.io/progrium/linux-build:latest AS kernel
COPY kernel.config .config
RUN make ARCH=i386 CROSS_COMPILE=i686-linux-gnu- oldconfig < /dev/null && \
    make ARCH=i386 CROSS_COMPILE=i686-linux-gnu- bzImage -j$(nproc)

FROM ghcr.io/progrium/v86:latest AS v86

FROM scratch AS bundle
COPY --from=rootfs / /rootfs
COPY --from=kernel /build/arch/x86/boot/bzImage /kernel/bzImage
COPY --from=v86 /libv86.js /v86/libv86.js
COPY --from=v86 /v86.wasm /v86/v86.wasm
COPY --from=v86 /bios/seabios.bin /v86/seabios.bin
COPY --from=v86 /bios/vgabios.bin /v86/vgabios.bin
COPY ./init.js /init.js