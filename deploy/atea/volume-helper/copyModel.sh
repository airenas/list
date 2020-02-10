#!/bin/bash
###########################################################################################
# Helper script to upload one model to kubernetes enviromnment
###########################################################################################
# first: run the script to volume helper container
# : kubectl apply -f helper.yml
# after: run the script to destroy volume helper container
# : kubectl delete deployment vh
###########################################################################################
localDir=$1
remoteDir=$2
###########################################################################################
podName=$(kubectl get po | grep -e '^vh' | head -n 1 | awk '{print $1}')
echo "Pod name = $podName"
if [ "$podName" == "" ] ; then
   echo -e "No pod found!!!\nDid you run: \n\nkubectl apply -f helper.yml ?\n"
   exit 1
fi
echo "From     = $localDir"
if [ "$localDir" == "" ] ; then
   echo -e "No local dir provided! Usage: copyModel.sh <localDir> <remoteDir>"
   exit 1
fi
echo "To       = $remoteDir"
if [ "$2" == "" ] ; then
   echo -e "No local dir provided! Usage: copyModel.sh <localDir> <remoteDir>"
   exit 1
fi

###########################################################################################
echo -e "\n\ncreating remote dir = $remoteDir"
kubectl exec $podName -i -- mkdir -p $remoteDir
###########################################################################################
echo -e "\n\ncopy kaldi models = $localDir"
rsync -avurP --blocking-io --rsync-path=$remoteDir --rsh="kubectl exec $podName -i -- " $localDir rsync:$remoteDir
###########################################################################################
echo -e "\n\nDone.\n\nNow unload helper container:\nkubectl delete deployment vh\n"
###########################################################################################
