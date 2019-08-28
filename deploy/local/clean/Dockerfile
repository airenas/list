FROM iron/base
LABEL Description="This image is used to start the LiST data cleaning service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

VOLUME ["/fileStorage"]
WORKDIR /app

COPY deploy/cleanService /app/
COPY config.yaml /app/

ENTRYPOINT ["./cleanService"]
