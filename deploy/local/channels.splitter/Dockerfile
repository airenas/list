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

COPY deploy/send.metric /app/
COPY deploy/*.sh /app/

COPY config.yaml /app/
COPY deploy/cmd/Makefile /cmd/
COPY Makefile.options /cmd/

RUN chown app:app /app/* /app
RUN chown app:app /cmd/* /cmd

USER app

ENTRYPOINT ["./cmdWorkerService"]
