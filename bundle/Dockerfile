FROM i386/alpine:latest
RUN apk add -u e2fsprogs podman netavark
COPY ./bin/init /bin/init
COPY ./bin/post-dhcp /bin/post-dhcp
COPY ./etc/containers.conf /etc/containers/containers.conf