#!/bin/bash
###############################################################
### Extract result for transcription ID
###############################################################
### takes the ID and wav file as input
### the result is put into <input>.txt file
###############################################################
### change the url to the correct one
url=https://prn509.vdu.lt:7080/ausis
###############################################################
id=$1
file=$2
###############################################################
resultURL="$url/result.service/result"
#resultURL=""http://127.0.0.1:8003/result""
###############################################################
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

###############################################################
echo "File: $file"
echo "Getting result..."
curl -X GET -k $resultURL/$id/result.txt -o "$file.txt" 2>/dev/null
if [ $? -gt "0" ] ; then
   echo -e "${RED}FAILED $file\n\tCan't get file.${NC}"
   exit 1
else
   echo -e "${GREEN}DONE $file${NC}"
fi
###############################################################
