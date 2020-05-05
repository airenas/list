FROM airenas/ubuntu.make:2.0
LABEL Description="This image is used to host the LiST result prepare worker service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="2.0" 

VOLUME ["/data"]
VOLUME ["/apps"]
VOLUME ["/models"]
WORKDIR /app

COPY deploy/cmdWorkerService /app/
COPY deploy/punct.lattice /app/
COPY deploy/lattice.to.text /app/
COPY deploy/lattice.to.webvtt /app/
COPY deploy/send.metric /app/
COPY config.yaml /app/
COPY deploy/*.sh /app/

COPY deploy/cmd/Makefile /cmd/
COPY Makefile.options /cmd/
COPY deploy/cmd/restore /cmd/restore


ENTRYPOINT ["./cmdWorkerService"]
