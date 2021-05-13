FROM airenas/ubuntu.make:2.0 as worker

RUN apt-get update && apt-get install -y sox ffmpeg bc

WORKDIR /app
COPY deploy/cmdWorkerService /app/
COPY deploy/send.metric /app/
COPY deploy/*.sh /app/

FROM worker

LABEL Description="This image is used to host the LiST audio merger service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 
VOLUME ["/data"]
VOLUME ["/models"]

COPY config.yaml /app/
COPY deploy/cmd/Makefile /cmd/
COPY Makefile.options /cmd/

ENTRYPOINT ["./cmdWorkerService"]
