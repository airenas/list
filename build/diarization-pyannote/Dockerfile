#########################################################################################
## docker will invoke this file from ../.. dir in order to access code
#########################################################################################
FROM golang:1.19-alpine as builder

ARG TOOLS_VERSION=0.1

WORKDIR /go/src/
ENV CGO_ENABLED=0

COPY ./src/tools /go/src

RUN go build -o /go/bin/send.metric -ldflags "-X main.version=$TOOLS_VERSION" send.metric/main.go
RUN go build -o /go/bin/fix.segments -ldflags "-X main.version=$TOOLS_VERSION" fix.segments/main.go
#########################################################################################
FROM airenas/list-worker:1.0.400-44fd44b as cmdRunner
FROM kaldiasr/kaldi:cpu-debian10 as kaldi
#########################################################################################
FROM python:3.10-bookworm as runner

ARG BUILD_VERSION=0.1

LABEL org.opencontainers.image.version=$BUILD_VERSION \
      org.opencontainers.image.authors="airenass@gmail.com" \
      name="result-make-service" \
      maintainer="airenass@gmail.com" \
      vendor="airenass@gmail.com" \
      version=$BUILD_VERSION \
      release=$BUILD_VERSION \
      summary="This image is used to start the LiST diarization-pyannote worker service" \
      description="This image is used to start the LiST diarization-pyannote worker service" 

RUN apt-get update && apt-get install -y \
      build-essential \
      libsndfile1



VOLUME ["/data"]
VOLUME ["/models"]
WORKDIR /app

COPY src/decode/diarization-pyannote/python/requirements.txt /cmd/python/
RUN pip install -r /cmd/python/requirements.txt

COPY --from=cmdRunner /app/cmdWorkerService /app/
COPY --from=kaldi //opt/kaldi/egs/wsj/s5/utils/utt2spk_to_spk2utt.pl /app/
COPY --from=builder /go/bin/send.metric /app/
COPY --from=builder /go/bin/fix.segments /app/
COPY build/diarization-pyannote/config.yaml /app/
COPY src/decode/scripts/*sh /app/
COPY src/decode/diarization-pyannote/python/scripts/*.py /app/scripts/

COPY build/diarization-pyannote/Makefile.options /cmd/
COPY src/decode/diarization-pyannote/Makefile /cmd/

ENTRYPOINT ["./cmdWorkerService"]
