#!/bin/bash
###########################################################################################
# Helper script to upload models to kubernetes enviromnment
###########################################################################################
# first: run the script to volume helper container
# : kubectl apply -f helper.yml
# after: run the script to destroy volume helper container
# : kubectl delete deployment vh
###########################################################################################
# slashes at the end required!
diarizationModels=diarization_models/
kaldiModels=kaldi_models/
punctuationModels=punctuation_models/
apps=kaldi_apps/
###########################################################################################
podName=$(kubectl get po | grep -e '^vh' | head -n 1 | awk '{print $1}')
echo "Pod name = $podName"
if [ "$podName" == "" ] ; then
   echo -e "No pod found!!!\nDid you run: \n\nkubectl apply -f helper.yml ?\n"
   exit 1
fi

###########################################################################################
echo -e "\n\ncopy diarization models = $diarizationModels"
rsync -avurP --blocking-io --rsync-path=/dmodels --rsh="kubectl exec $podName -i -- " $diarizationModels rsync:/dmodels
###########################################################################################
echo -e "\n\ncopy apps = $apps"
rsync -avurP --blocking-io --rsync-path=/apps --rsh="kubectl exec $podName -i -- " $apps rsync:/apps
###########################################################################################
echo -e "\n\ncopy kaldi models = $kaldiModels"
rsync -avurP --blocking-io --rsync-path=/models --rsh="kubectl exec $podName -i -- " $kaldiModels rsync:/models
###########################################################################################
echo -e "\n\ncopy punctuation models = $punctuationModels"
rsync -avurP --blocking-io --rsync-path=/pmodels --rsh="kubectl exec $podName -i -- " $punctuationModels rsync:/pmodels
###########################################################################################
echo -e "\n\nDone.\n\nNow unload helper container:\nkubectl delete deployment vh\n"
###########################################################################################
