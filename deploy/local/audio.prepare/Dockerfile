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
FROM airenas/list-worker:1.0.400-44fd44b as cmdRunner

FROM airenas/ubuntu.make:2.0 as runner

RUN apt-get update && apt-get install -y sox ffmpeg bc

WORKDIR /app
COPY --from=cmdRunner /app/cmdWorkerService /app/
COPY --from=builder /go/bin/send.metric /app/
COPY deploy/local/audio.prepare/deploy/*.sh /app/

LABEL Description="This image is used to host the LiST audio prepare worker service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="2.0" 

VOLUME ["/data"]
VOLUME ["/models"]

COPY deploy/local/audio.prepare/deploy/sfg /app/
COPY deploy/local/audio.prepare/config.yaml /app/
COPY deploy/local/audio.prepare/deploy/cmd/Makefile /cmd/
COPY deploy/local/audio.prepare/Makefile.options /cmd/

ENTRYPOINT ["./cmdWorkerService"]
