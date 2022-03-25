FROM iron/base:3.14.0

ARG BUILD_VERSION=0.1
LABEL org.opencontainers.image.version=$BUILD_VERSION \
      org.opencontainers.image.authors="airenass@gmail.com" \
      name="list-upload" \
      maintainer="airenass@gmail.com" \
      vendor="airenass@gmail.com" \
      version=$BUILD_VERSION \
      release=$BUILD_VERSION \
      summary="This image is used to start the LiST upload service" \
      description="This image is used to start the LiST upload service" 

ENV UID=999
VOLUME /data
VOLUME /models
EXPOSE 8000

WORKDIR /app

RUN addgroup app && adduser -S -G app --uid $UID app
RUN mkdir -p /app && chown -R app:app /app   

COPY deploy/uploadService /app/
COPY config.yaml /app/

RUN chown app:app /app/* /app
USER app

ENTRYPOINT ["./uploadService"]
