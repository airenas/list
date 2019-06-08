FROM ubuntu:16.04

LABEL Description="This image is used to start the LiST kafka integration service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY deploy/kafkaIntegrationService /app/
COPY config.yaml /app/

ENTRYPOINT ["./kafkaIntegrationService"]
