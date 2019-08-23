FROM iron/base
LABEL Description="This image is used to start the LiST punctuation service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

VOLUME ["/model"]
WORKDIR /app

COPY deploy/punctuationService /app/
COPY config.yaml /app/

ENTRYPOINT ["./punctuationService"]
