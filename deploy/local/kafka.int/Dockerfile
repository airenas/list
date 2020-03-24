FROM ubuntu:18.04

LABEL Description="This image is used to start the LiST kafka integration service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" 

WORKDIR /app

COPY deploy/kafkaIntegrationService /app/
COPY config.yaml /app/

ENTRYPOINT ["./kafkaIntegrationService"]
