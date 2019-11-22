FROM iron/base

LABEL Description="This image is used to start the LiST upload service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

VOLUME ["/data"]
VOLUME ["/models"]
WORKDIR /app

COPY deploy/uploadService /app/
COPY config.yaml /app/

EXPOSE 8000

CMD ["./uploadService"]
