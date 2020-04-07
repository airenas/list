FROM airenas/ubuntu.make:2.0
LABEL Description="This image is used to host the LiST transcription worker service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="2.0" 

VOLUME ["/data"]
VOLUME ["/apps"]
VOLUME ["/models"]
WORKDIR /app

COPY deploy/cmdWorkerService /app/
COPY config.yaml /app/
COPY deploy/pipe.runner /app/

COPY deploy/cmd/Makefile /cmd/
COPY Makefile.options /cmd/

ENTRYPOINT ["./cmdWorkerService"]
