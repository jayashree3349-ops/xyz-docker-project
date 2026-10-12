FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY service-a/main.go ./service-a/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/service-a ./service-a

FROM alpine:3.22
RUN adduser -D -u 10001 appuser
COPY --from=builder /out/service-a /usr/local/bin/service
USER 10001
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/service"]
