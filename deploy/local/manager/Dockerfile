FROM iron/base

LABEL Description="This image is used to start the LiST transcription manager service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY deploy/managerService /app/
COPY config.yaml /app/

#CMD ls -la
#ENTRYPOINT ["ls", "-la"]
ENTRYPOINT ["./managerService"]
