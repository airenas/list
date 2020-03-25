FROM iron/base

LABEL Description="This image is used to start the LiST upload service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas"

VOLUME ["/data"]
VOLUME ["/models"]
WORKDIR /app

COPY deploy/uploadService /app/
COPY config.yaml /app/

EXPOSE 8000

CMD ["./uploadService"]
