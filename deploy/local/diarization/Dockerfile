FROM ubuntu:16.04
LABEL Description="This image is used to host the LiST diarization worker service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

VOLUME ["/data"]
VOLUME ["/apps"]
VOLUME ["/models"]
WORKDIR /app

COPY deploy/cmdWorkerService /app/
COPY config.yaml /app/

COPY deploy/cmd/Makefile /cmd/
COPY Makefile.options /cmd/

RUN apt-get update && apt-get install -y build-essential sox

RUN DEBIAN_FRONTEND=noninteractive \ 
      apt-get -y install default-jre-headless && \ 
      apt-get clean && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["./cmdWorkerService"]
