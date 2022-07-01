#########################################################################################
## docker will invoke this file from ../.. dir in order to access code
#########################################################################################
FROM golang:1.17-alpine as builder

ARG TOOLS_VERSION=0.1

WORKDIR /go/src/
ENV CGO_ENABLED=0

COPY ./src/tools /go/src

RUN go build -o /go/bin/send.metric -ldflags "-X main.version=$TOOLS_VERSION" send.metric/main.go
RUN go build -o /go/bin/punct.lattice -ldflags "-X main.version=$TOOLS_VERSION" punct.lattice/main.go
RUN go build -o /go/bin/lattice.to.text -ldflags "-X main.version=$TOOLS_VERSION" lattice.to.text/main.go
RUN go build -o /go/bin/empty.lattice -ldflags "-X main.version=$TOOLS_VERSION" empty.lattice/main.go
RUN go build -o /go/bin/fix.lattice.time -ldflags "-X main.version=$TOOLS_VERSION" fix.lattice.time/main.go
RUN go build -o /go/bin/lattice.to.webvtt -ldflags "-X main.version=$TOOLS_VERSION" lattice.to.webvtt/main.go
RUN go build -o /go/bin/int.to.word -ldflags "-X main.version=$TOOLS_VERSION" int.to.word/main.go
#########################################################################################
FROM airenas/list-worker:1.0.400-44fd44b as cmdRunner

#########################################################################################
FROM airenas/ubuntu.make:2.0 as runner

ARG BUILD_VERSION=0.1

LABEL org.opencontainers.image.version=$BUILD_VERSION \
      org.opencontainers.image.authors="airenass@gmail.com" \
      name="result-make-service" \
      maintainer="airenass@gmail.com" \
      vendor="airenass@gmail.com" \
      version=$BUILD_VERSION \
      release=$BUILD_VERSION \
      summary="This image is used to start the LiST result prepare worker service" \
      description="This image is used to start the LiST result prepare worker service" 

VOLUME ["/data"]
VOLUME ["/apps"]
VOLUME ["/models"]
WORKDIR /app

COPY --from=cmdRunner /app/cmdWorkerService /app/
COPY --from=builder /go/bin/send.metric /app/
COPY --from=builder /go/bin/punct.lattice /app/
COPY --from=builder /go/bin/lattice.to.text /app/
COPY --from=builder /go/bin/empty.lattice /app/
COPY --from=builder /go/bin/fix.lattice.time /app/
COPY --from=builder /go/bin/lattice.to.webvtt /app/
COPY --from=builder /go/bin/int.to.word /app/
COPY deploy/local/result.make/config.yaml /app/
COPY deploy/local/result.make/deploy/*.sh /app/

COPY deploy/local/result.make/deploy/cmd/Makefile /cmd/
COPY deploy/local/result.make/Makefile.options /cmd/
COPY deploy/local/result.make/deploy/cmd/restore /cmd/restore

ENTRYPOINT ["./cmdWorkerService"]
