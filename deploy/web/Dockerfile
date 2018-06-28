FROM nginx:latest
MAINTAINER Airenas Vaiciunas <airenass@gmail.com>

LABEL Description="This image is used to start the ear-app web service" Vendor="Airenas Vaiciunas" Version="1.0"

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY deploy/html /usr/share/nginx/html
