#!/bin/bash
###############################################################
### Send one file to transcription
###############################################################
### takes the wav file as input
### passes it to the transcriber, 
### the result is <transcription id> <input file>
###############################################################
### change the url to the correct one
url=http://localhost:8001
###############################################################
file=$1
echoerr() { echo "$@" 1>&2; }
###############################################################
uploadURL="$url/upload"
###############################################################
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color
echoerr() { echo "$@" 1>&2; }
###############################################################
echoerr "File: $file"
echoerr "Uploading..."
id=$(curl -X POST -k $uploadURL -H 'Accept: application/json' -H 'content-type: multipart/form-data' -F file=@$file 2>/dev/null | jq -r '.["id"]')
if [ $? -gt "0" ] ; then
   echoerr -e "${RED}FAILED $file\n\tCan't send file.${NC}"
   exit 1
fi
echo "$id $file"
###############################################################
