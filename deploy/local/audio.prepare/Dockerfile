FROM airenas/ubuntu.make:2.0
LABEL Description="This image is used to host the LiST audio prepare worker service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="2.0" 

VOLUME ["/data"]
VOLUME ["/models"]
WORKDIR /app

COPY deploy/cmdWorkerService /app/
COPY config.yaml /app/

COPY deploy/cmd/Makefile /cmd/
COPY Makefile.options /cmd/

RUN apt-get update && apt-get install -y sox ffmpeg bc

ENTRYPOINT ["./cmdWorkerService"]
