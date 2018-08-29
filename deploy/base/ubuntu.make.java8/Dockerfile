FROM airenas/ubuntu.make:1.0
LABEL Description="This image is used as base with Ubuntu, Make and Java 8 installed" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

RUN apt-get update && apt-get install -y build-essential sox

RUN DEBIAN_FRONTEND=noninteractive \ 
      apt-get -y install default-jre-headless && \ 
      apt-get clean && rm -rf /var/lib/apt/lists/*