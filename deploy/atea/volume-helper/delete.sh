#!/bin/bash
###########################################################################################
# Helper script to delete all data from kubernetes enviromnment
###########################################################################################
# first: run the script to volume helper container
# : kubectl apply -f helper.yml
# after: run the script to destroy volume helper container
# : kubectl delete deployment vh
###########################################################################################
# no slash at the end!
###########################################################################################
podName=$(kubectl get po | grep -e '^vh' | head -n 1 | awk '{print $1}')
echo "Pod name = $podName"
if [ "$podName" == "" ] ; then
   echo -e "No pod found!!!\nDid you run: \n\nkubectl apply -f helper.yml ?\n"
   exit 1
fi
###########################################################################################
read -p "Delete all transcription data/model files [no/yes]: " confirm
confirm=${confirm:-no}
if [ "$confirm" != "yes" ] ; then
   echo -e "Skip deletion!\nYou pressed $confirm\n"
   exit 1
fi

kubectl exec $podName -i -- rm -r /apps
kubectl exec $podName -i -- rm -rf /rabbitmq
kubectl exec $podName -i -- rm -rf /mongo
kubectl exec $podName -i -- rm -rf /models
kubectl exec $podName -i -- rm -rf /dmodels
kubectl exec $podName -i -- rm -rf /filestorage

###########################################################################################
echo -e "\n\nDone.\n\nNow unload helper container:\nkubectl delete deployment vh\n"
###########################################################################################
