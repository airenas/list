#!/bin/bash
###############################################################
### Extract result for transcription ID
###############################################################
### takes the ID and wav file as input
### the result is put into <input>.txt file
###############################################################
### change the url to the correct one
url=$RESULT_SERVICE
###############################################################
result_dir=$1
id=$2
file=${3##*/}
file=${file%.*}.txt
###############################################################
resultURL="$url/result"
###############################################################
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

###############################################################
echo "File: $file"
echo "Getting result..."
code=$(curl -X GET -k $resultURL/$id/result.txt -o "$result_dir/$file" 2>/dev/null -w '%{http_code}')
res=$?
echo $code
if [ "$code" == "404" ] || [ "$code" == "000" ] || [ $res -gt "0" ]; then
   echo -e "${RED}FAILED $file\t\tCan't get file.${NC}"
   rm -f "$result_dir/$file"
   exit 1
fi
echo -e "${GREEN}DONE $file${NC}"
###############################################################
