FROM ubuntu:16.04
LABEL Description="This image is used to host the LiST transcription worker service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

VOLUME ["/data"]
VOLUME ["/apps"]
VOLUME ["/models"]
WORKDIR /app

COPY deploy/cmdWorkerService /app/
COPY config.yaml /app/

COPY deploy/cmd/Makefile /cmd/
COPY Makefile.options /cmd/

RUN apt-get update && apt-get install -y build-essential libatlas-base-dev

ENTRYPOINT ["./cmdWorkerService"]
