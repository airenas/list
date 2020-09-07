FROM semantikadocker.vdu.lt/diarization-alt:0.2
LABEL Description="This image is used to host the LiST diarization-alt worker service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="0.1" 

VOLUME ["/data"]
VOLUME ["/models"]
WORKDIR /app

RUN apt-get update && apt-get install -y build-essential

COPY deploy/cmdWorkerService /app-list/
COPY config.yaml /app-list/
COPY deploy/send.metric /app-list/
COPY deploy/fix.segments /app-list/
COPY deploy/*.sh /app-list/

COPY deploy/cmd/Makefile /cmd/
COPY Makefile.options /cmd/

ENTRYPOINT ["/app-list/cmdWorkerService"]
