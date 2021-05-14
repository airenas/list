FROM airenas/ubuntu.make:2.0
LABEL Description="This image is used to host the LiST zoom join result worker service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

VOLUME ["/data"]
WORKDIR /app

COPY deploy/cmdWorkerService /app/
COPY deploy/send.metric /app/
COPY config.yaml /app/
COPY deploy/*.sh /app/

COPY deploy/cmd/Makefile /cmd/
COPY Makefile.options /cmd/

ENTRYPOINT ["./cmdWorkerService"]
