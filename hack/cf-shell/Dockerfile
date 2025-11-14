FROM golang:1.25.0-alpine AS build
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /cfshell .

FROM alpine:latest
RUN apk add --no-cache bash coreutils tini fuse
ENV GO_VERSION="1.25.1"
WORKDIR /usr/local
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -xzf go${GO_VERSION}.linux-amd64.tar.gz go/src go/pkg go/bin \
    && rm go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH=/usr/local/go/bin:$PATH
RUN go telemetry off && go build std
COPY --from=build /cfshell /bin/cfshell
EXPOSE 8080
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/bin/cfshell"]