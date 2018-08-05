FROM iron/base
MAINTAINER Airenas Vaiciunas <airenass@gmail.com>

LABEL Description="This image is used to start the LiST status provider service" Vendor="Airenas Vaiciunas" Version="1.0"

WORKDIR /app

COPY deploy/statusProviderService /app/
COPY config.yaml /app/

EXPOSE 8000

#CMD ls -la
#ENTRYPOINT ["ls", "-la"]
ENTRYPOINT ["./statusProviderService"]
