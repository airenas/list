# Deploy information of AFT components for ATEA Kubernetes

The deployment includes:
- k8s configurations for volumes, rabbitmq, mongo db, aft transcription services and integration with kafka
- scripts for uploading model files

---
## Prerequisite
- *kubectl* must be installed and configured to connect to ATEA k8s. For example `kubectl get pods` works without errors
- all sample scripts in this document assume you are in the current directory

## Deploy instructions
###Secrets
1. Prepare secrets. 
Go to *secrets* directory and prepare *secrets.yml* file from *secrets.yml.template* replacing  *{{rabbit-mq-user}}*, *{{rabbit-mq-pass}}*, *{{mongo-user}}*, *{{mongo-pass}}*. Or you can use sample scripts to generate random passwords: `(cd secrets && ./generate.sh)`. 
**Note:** values must be  base64 encoded.
2. Install secrets to k8s:
```bash
kubectl apply -f secrets
```
3. Delete *secrets.yml* file:
```bash
rm secrets/secrets.yml
```
###Volumes
The system requires kaldi binaries and models files to be populated into k8s volumes. Assume you have data locally and you can access *kaldi apps*, *kaldi models* and *diarization models* directories. 
1. Init volumes: `kubectl apply -f deployments/volumes.yml`
2. Run temporary pod which attaches to volumes: `kubectl apply -f volume-helper`
3. Configure script *volume-helper/copy.sh*, populate correct local directories for *diarizationModels*, *kaldiModels* and *apps*.
4. Copy everything to k8s volumes: 
```bash
(cd volume-helper && ./copy.sh)
```
5. Remove temporary volume helper pod: `kubectl delete deployment vh`

###Services
Deploy services: `kubectl apply -f deployments`

Procedure for update/redeploy services: 
1. Make required changes, for example update service versions in *deployments/transcription.yml*
2. Commit changes to git 
3. Run `kubectl apply -f deployments`

##Additional info



###Sample scripts for troubleshooting
**Service not starting:** `kubectl describe deployment <deployment-name>`
**Get log for some pod:** `kubectl logs pod <pod-name>`
**Attach to running pod:** `kubectl exec -it <pod-name> /bin/bash`. Some container are very slim and do not include *bash*, try: `kubectl exec -it <pod-name> /bin/sh`. 

###Access services through kubectl port forward
To access services **upload**, **status**, **result** running on k8s through local ports 8001, 8002, 8003:
- `kubectl -n aft port-forward svc/upload 8001:8000`
- `kubectl -n aft port-forward svc/status 8002:8000`
- `kubectl -n aft port-forward svc/result 8003:8000`

**Airenas Vaičiūnas**

* [bitbucket.org/airenas](https://bitbucket.org/airenas)
* [linkedin.com/in/airenas](https://www.linkedin.com/in/airenas/)

---