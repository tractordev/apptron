ARG BUILDPLATFORM
ARG TARGETPLATFORM

FROM --platform=$BUILDPLATFORM ghcr.io/progrium/linux-build:latest AS build
COPY kernel.config .config
RUN make ARCH=i386 CROSS_COMPILE=i686-linux-gnu- oldconfig < /dev/null && \
    make ARCH=i386 CROSS_COMPILE=i686-linux-gnu- bzImage -j$(nproc)

FROM scratch AS kernel
COPY --from=build /build/arch/x86/boot/bzImage /bzImage
CMD ["true"]