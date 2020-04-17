FROM iron/base

LABEL Description="This image is used to start the LiST queue dispather service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

VOLUME ["/data"]
VOLUME ["/models"]

WORKDIR /app

COPY deploy/dispatcherService /app/
COPY config.yaml /app/

ENTRYPOINT ["./dispatcherService"]
