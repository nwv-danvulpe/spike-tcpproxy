FROM golang:1.15-alpine as builder

WORKDIR /workspace

COPY . .

RUN go build -o spike-tcpproxy && \
    chmod +x spike-tcpproxy

FROM alpine:latest

COPY --from=builder /workspace/spike-tcpproxy /usr/bin/spike-tcpproxy

RUN apk --no-cache add curl

ENTRYPOINT [ "/usr/bin/spike-tcpproxy" ]