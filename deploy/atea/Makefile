-include Makefile.options
srv?=kafka
nth?=1
#####################################################################################
LOCAL_PORT_UPLOAD?=8001
LOCAL_PORT_STATUS?=8002
LOCAL_PORT_RESULT?=8003
LOCAL_PORT_PROMETHEUS?=8004
LOCAL_PORT_RABBITMQ?=8005
#####################################################################################
show-pods: 
	kubectl get pods
show-deployments: 
	kubectl get deployments
show-services: 
	kubectl get services
#####################################################################################
install-vh: 
	kubectl apply -f volume-helper
clean-vh: 
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	kubectl delete -f volume-helper
#####################################################################################
proxy-upload:
	kubectl -n aft port-forward svc/upload $(LOCAL_PORT_UPLOAD):8000
info-upload:
	curl -i http://localhost:$(LOCAL_PORT_UPLOAD)/live?full=1
test-upload:
	./utils/test.sh proxy_upload info_upload
#####################################################################################
proxy-status:
	kubectl -n aft port-forward svc/status $(LOCAL_PORT_STATUS):8000
info-status:
	curl -i http://localhost:$(LOCAL_PORT_STATUS)/live?full=1
test-status:
	./utils/test.sh proxy_status info_status	
#####################################################################################
proxy-result:
	kubectl -n aft port-forward svc/result $(LOCAL_PORT_RESULT):8000
info-result:
	curl -i http://localhost:$(LOCAL_PORT_RESULT)/live?full=1
test-result:
	./utils/test.sh proxy_result info_result		
#####################################################################################
proxy-prometheus:
	kubectl -n aft port-forward svc/prometheus $(LOCAL_PORT_PROMETHEUS):9090
proxy-rabbitmq:
	kubectl -n aft port-forward svc/rabbitmq $(LOCAL_PORT_RABBITMQ):15672
#####################################################################################
vh_pod_name:
	$(eval vh_pod_name := $(shell kubectl get po | grep -e '^vh' | head -n 1 | awk '{print $$1}'))
	@$(if $(strip $(vh_pod_name)),echo "Pod name = $(vh_pod_name)",echo No vh pod && exit 1)
vh_cmd: vh_pod_name
	$(eval vh_cmd := kubectl exec -it $(vh_pod_name))
#####################################################################################
connect-vh: vh_cmd
	$(vh_cmd) /bin/sh
#####################################################################################
connect: $(srv)_cmd
	echo "Connecting $($(srv)_pod_name)"
	$($(srv)_cmd) /bin/bash	
logs: $(srv)_pod_name
	echo "Connecting $($(srv)_pod_name)"
	kubectl logs $($(srv)_pod_name)
$(srv)_pod_name:
	$(eval $(srv)_pod_name := $(shell kubectl get po | grep -e '^$(srv)' | head -n $(nth) | tail -n 1 | awk '{print $$1}'))
	@$(if $(strip $($(srv)_pod_name)),echo "Pod name = $($(srv)_pod_name)",echo No $(srv) pod && exit 1)
$(srv)_cmd: $(srv)_pod_name
	$(eval $(srv)_cmd := kubectl exec -it $($(srv)_pod_name))	
#####################################################################################
copy_models: 
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/ac/20_04_17-2g /models/ac/
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/ac/19_11_18-3g /models/ac/
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/ac/graph_w63g /models/ac/

copy_rnnlm: 
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/rnnlm/v1 /models/rnnlm/

copy_lm: 
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/lm/v3 /models/lm/
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/lm/v1 /models/lm/

#####################################################################################
prepare-secrets:
	cd secrets && ./generate.sh
install-secrets:
	kubectl apply -f secrets
clean-secrets:
	rm -f secrets/secrets.yml
#####################################################################################
install-volumes: install-vh
	kubectl apply -f deployments/volumes.yml
########### DATA ####################################################################
copy-apps: vh_cmd
	$(vh_cmd) make apps=apps-1.0 install-apps
init-storage: vh_cmd
	$(vh_cmd) make init-storage
copy-config:
	./volume-helper/copyModel.sh models/config/ /models/config
copy-data: copy-apps copy-config init-storage copy-models copy_lm copy-diarization copy-punctuation
clean-data: vh_cmd
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	$(vh_cmd) make clean-data clean-apps
#####################################################################################
install-services:
	kubectl apply -f deployments/
clean-services:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	kubectl delete -f deployments/	
#####################################################################################
delete-prometheus:
	kubectl delete deployment prometheus | true
restart-prometheus: delete-prometheus
	kubectl apply -f deployments/metrics.yml
#####################################################################################
.PHONY: show-pods install-vh clean-vh 
