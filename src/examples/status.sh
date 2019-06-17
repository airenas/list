#!/bin/bash
###############################################################
### Get transcription status
###############################################################
### takes the ID and wav file as input
### returns status
###############################################################
### change the url to the correct one
url=https://prn509.vdu.lt:7080/ausis
###############################################################
id=$1
file=$2
###############################################################
statusURL="$url/status.service/status"
#statusURL="http://127.0.0.1:8002/status"
###############################################################
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color
###############################################################
r=$(curl -X GET -k $statusURL/$id -H "accept: application/json" 2>/dev/null)
err=$(echo "$r" | jq -r '.["error"]')
status=$(echo $r | jq -r '.["status"]')
if [ -n "$err" ] ; then
   echo -e "${RED}FAILED $ID\n\t$err${NC}"
fi
if [ "$status" == "COMPLETED" ] ; then
   echo -e "$id ${GREEN}${status}${NC}"
else
   echo "$id $status"
fi
