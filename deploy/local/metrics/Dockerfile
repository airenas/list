FROM iron/base
LABEL Description="This image is used to host the LiST metrics collector service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY deploy/metricsCollector /app/
COPY config.yaml /app/

ENTRYPOINT ["./metricsCollector"]
