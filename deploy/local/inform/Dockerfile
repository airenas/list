FROM iron/base:3.14

ARG BUILD_VERSION=0.1
LABEL org.opencontainers.image.version=$BUILD_VERSION \
      org.opencontainers.image.authors="airenass@gmail.com" \
      name="list-upload" \
      maintainer="airenass@gmail.com" \
      vendor="airenass@gmail.com" \
      version=$BUILD_VERSION \
      release=$BUILD_VERSION \
      summary="This image is used to start the LiST inform service" \
      description="This image is used to start the LiST inform service" 

RUN apk --no-cache add tzdata

ENV UID=999
RUN addgroup app && adduser -S -G app --uid $UID app
RUN mkdir -p /app && chown -R app:app /app   

WORKDIR /app

COPY deploy/informService /app/
COPY config.yaml /app/
COPY mail.tmpl /app/template/

RUN chown app:app /app/* /app
USER app

ENTRYPOINT ["./informService"]
