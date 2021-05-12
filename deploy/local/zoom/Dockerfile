FROM iron/base

LABEL Description="This image is used to start the LiST zoom manager service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY deploy/zoomService /app/
COPY config.yaml /app/

ENTRYPOINT ["./zoomService"]
