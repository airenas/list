#########################################################################################
## docker will invoke this file from ../.. dir in order to access code
#########################################################################################
FROM golang:1.17-alpine as builder

ARG TOOLS_VERSION=0.1

WORKDIR /go/src/
ENV CGO_ENABLED=0

COPY ./src/tools /go/src

RUN go build -o /go/bin/send.metric -ldflags "-X main.version=$TOOLS_VERSION" send.metric/main.go
#########################################################################################
FROM airenas/list-worker:1.0.400-44fd44b as runner

USER root

RUN apk add --no-cache ffmpeg

VOLUME ["/data"]
VOLUME ["/models"]

ARG BUILD_VERSION=0.1

LABEL org.opencontainers.image.version=$BUILD_VERSION \
      org.opencontainers.image.authors="airenass@gmail.com" \
      name="channels-splitter" \
      maintainer="airenass@gmail.com" \
      vendor="airenass@gmail.com" \
      version=$BUILD_VERSION \
      release=$BUILD_VERSION \
      summary="This image is used to start the LiST channels splitter worker" \
      description="This image is used to start the LiST channels splitter worker" 

COPY --from=builder /go/bin/send.metric /app/
COPY deploy/local/channels.splitter/deploy/*.sh /app/

COPY deploy/local/channels.splitter/config.yaml /app/
COPY deploy/local/channels.splitter/deploy/cmd/Makefile /cmd/
COPY deploy/local/channels.splitter/Makefile.options /cmd/

RUN chown app:app /app/* /app
RUN chown app:app /cmd/* /cmd

USER app

ENTRYPOINT ["./cmdWorkerService"]
