FROM iron/base
MAINTAINER Airenas Vaiciunas <airenass@gmail.com>

LABEL Description="This image is used to start the LiST upload service" Vendor="Airenas Vaiciunas" Version="1.0"

VOLUME ["/data"]
WORKDIR /app

COPY deploy/uploadService /app/
COPY config.yaml /app/

EXPOSE 8000

#CMD ls -la
#ENTRYPOINT ["ls", "-la"]
ENTRYPOINT ["./uploadService"]
