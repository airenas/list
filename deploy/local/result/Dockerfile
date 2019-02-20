FROM iron/base
LABEL Description="This image is used to start the LiST result service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY deploy/resultService /app/
COPY config.yaml /app/

EXPOSE 8000

#CMD ls -la
#ENTRYPOINT ["ls", "-la"]
ENTRYPOINT ["./resultService"]
