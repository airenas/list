#!/bin/bash
echo Starting
sed -i "s|<base href=\"/\">|<base href=\"$BASE_HREF\">|" /usr/share/nginx/html/index.html
echo Updated base ref to $BASE_HREF
echo Starting nginx
nginx -g 'daemon off;'