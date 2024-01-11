FROM golang:1.21.4-alpine as builder
RUN apk -u add git
WORKDIR /src
COPY . ./
RUN ./build-server.sh

FROM alpine:latest
COPY --from=builder /src/server /server

ENTRYPOINT ["/server"]
