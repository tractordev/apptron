FROM golang:1.24.5-alpine AS build
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /cfnet .

FROM scratch
COPY --from=build /cfnet /cfnet
EXPOSE 8080
CMD ["/cfnet"]