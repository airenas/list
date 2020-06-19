FROM alpine:3.9.4

LABEL Description="This image is prepared for files copying to volumes with rsync" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

RUN apk update && apk add rsync wget make htop mc tar curl

VOLUME ["/utils"]
WORKDIR /utils

COPY utils/Makefile /utils/
COPY utils/Makefile.options /utils/


CMD ["tail", "-f", "/dev/null"]