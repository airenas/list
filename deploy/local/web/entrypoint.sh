#!/bin/bash
sed -i "s|<base href=\"/\">|<base href=\"$BASE_HREF\">|" /usr/share/nginx/html/index.html

nginx -g 'daemon off;'