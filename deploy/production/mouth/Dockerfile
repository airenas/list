FROM ubuntu:15.10
MAINTAINER Airenas Vaiciunas <airenass@gmail.com>
LABEL Description="This image is used to start the LiST membrane service" Vendor="Airenas Vaiciunas" Version="1.0"

VOLUME ["/data"]
WORKDIR /app

COPY deploy/mouth.service /app/
COPY config.yml /app/

EXPOSE 8000

CMD ./mouth.service