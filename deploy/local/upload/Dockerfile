FROM iron/base

LABEL Description="This image is used to start the LiST upload service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

VOLUME ["/data"]
WORKDIR /app

COPY deploy/uploadService /app/
COPY config.yaml /app/
COPY start.service.sh /app/

EXPOSE 8000

CMD ["./start.service.sh"]

#CMD ["./uploadService"]
