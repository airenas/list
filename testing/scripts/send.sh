#!/bin/bash
###############################################################
### Send one file to transcription
###############################################################
### takes the wav file as input
### passes it to the transcriber, 
### the result is <transcription id> <input file>
###############################################################
### change the url to the correct one
url=$UPLOAD_SERVICE
model=$1
###############################################################
file=$2
echoerr() { echo "$@" 1>&2; }
###############################################################
uploadURL="$url/upload"
###############################################################
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color
echoerr() { echo "$@" 1>&2; }
###############################################################
echoerr "File: $file model: $model url: $uploadURL"
echoerr "Uploading..."
numJoinParam='';if [ "${SKIP_NUM_JOIN}" == "1" ]; then numJoinParam='-F skipNumJoin=1'; fi
numSpeakersParam='';if [ -n "${NUMBER_OF_SPEAKERS}" ]; then numSpeakersParam="-F numberOfSpeakers=${NUMBER_OF_SPEAKERS}"; fi
echoerr "Params: ${numJoinParam} ${numSpeakersParam}"
c_res=$(curl -X POST -k $uploadURL -H 'Accept: application/json' -H 'content-type: multipart/form-data' -F recognizer=$model \
    -F file=@$file ${numJoinParam} ${numSpeakersParam} 2>/dev/null)
echoerr Response: $c_res
id=$(echo $c_res | jq -r '.["id"]')
if [ $? -gt "0" ] ; then
   echoerr -e "${RED}FAILED $file\n\tCan't send file.${NC}"
   echoerr -e "${RED}FAILED $id ${NC}"
   exit 1
fi
echo "$id $file"
###############################################################
