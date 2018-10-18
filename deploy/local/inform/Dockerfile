FROM iron/base

LABEL Description="This image is used to start the LiST inform service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY deploy/informService /app/
COPY config.yaml /app/

RUN apk --no-cache add tzdata

#CMD ls -la
#ENTRYPOINT ["ls", "-la"]
ENTRYPOINT ["./informService"]
