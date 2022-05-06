#########################################################################################
## docker will invoke this file from ../.. dir in order to access code
#########################################################################################
FROM golang:1.17-alpine as builder

ARG TOOLS_VERSION=0.1

WORKDIR /go/src/
ENV CGO_ENABLED=0

COPY ./src/tools /go/src

RUN go build -o /go/bin/send.metric -ldflags "-X main.version=$TOOLS_VERSION" send.metric/main.go
RUN go build -o /go/bin/lattices.to.webvtt -ldflags "-X main.version=$TOOLS_VERSION" lattices.to.webvtt/main.go
RUN go build -o /go/bin/lattices.to.text -ldflags "-X main.version=$TOOLS_VERSION" lattices.to.text/main.go
#########################################################################################
FROM airenas/list-worker:1.0.400-44fd44b as runner


VOLUME ["/data"]

ARG BUILD_VERSION=0.1

LABEL org.opencontainers.image.version=$BUILD_VERSION \
      org.opencontainers.image.authors="airenass@gmail.com" \
      name="join-results" \
      maintainer="airenass@gmail.com" \
      vendor="airenass@gmail.com" \
      version=$BUILD_VERSION \
      release=$BUILD_VERSION \
      summary="This image is used to start the LiST zoom join result worker" \
      description="This image is used to start the LiST zoom join result worker" 
      
USER root

COPY --from=builder /go/bin/send.metric /app/
COPY --from=builder /go/bin/lattices.to.webvtt /app/
COPY --from=builder /go/bin/lattices.to.text /app/
COPY deploy/local/join.results/deploy/*.sh /app/

COPY deploy/local/join.results/config.yaml /app/
COPY deploy/local/join.results/deploy/cmd/Makefile /cmd/
COPY deploy/local/join.results/Makefile.options /cmd/

RUN chown app:app /app/* /app
RUN chown app:app /cmd/* /cmd

USER app

ENTRYPOINT ["./cmdWorkerService"]
