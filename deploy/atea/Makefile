-include Makefile.options
#####################################################################################
LOCAL_PORT_UPLOAD?=8001
LOCAL_PORT_STATUS?=8002
LOCAL_PORT_RESULT?=8003
#####################################################################################
ku_po: 
	kubectl get pods

#####################################################################################
init_vh: 
	kubectl apply -f volume-helper
stop_vh: 
	kubectl delete deployment vh
#####################################################################################
copy_config:
	./volume-helper/copyModel.sh models/config/ /models/config
#####################################################################################
proxy_upload:
	kubectl -n aft port-forward svc/upload $(LOCAL_PORT_UPLOAD):8000
info_upload:
	curl -i http://localhost:$(LOCAL_PORT_UPLOAD)/live?full=1
test_upload:
	./utils/test.sh proxy_upload info_upload
#####################################################################################
proxy_status:
	kubectl -n aft port-forward svc/status $(LOCAL_PORT_STATUS):8000
info_status:
	curl -i http://localhost:$(LOCAL_PORT_STATUS)/live?full=1
test_status:
	./utils/test.sh proxy_status info_status	
#####################################################################################
proxy_result:
	kubectl -n aft port-forward svc/result $(LOCAL_PORT_RESULT):8000
info_result:
	curl -i http://localhost:$(LOCAL_PORT_RESULT)/live?full=1
test_result:
	./utils/test.sh proxy_result info_result		
#####################################################################################
vh_pod_name:
	$(eval vh_pod_name := $(shell kubectl get po | grep -e '^vh' | head -n 1 | awk '{print $$1}'))
	@$(if $(strip $(vh_pod_name)),echo "Pod name = $(vh_pod_name)",echo No vh pod && exit 1)
vh_cmd: vh_pod_name
	$(eval vh_cmd := kubectl exec -it $(vh_pod_name))
#####################################################################################
connect_vh: vh_cmd
	$(vh_cmd) /bin/sh
#####################################################################################
init_filestorage: vh_cmd
	$(vh_cmd) -- mkdir -p /filestorage/decoded/audio
	$(vh_cmd) -- mkdir -p /filestorage/decoded/trans
	$(vh_cmd) -- mkdir -p /filestorage/results
	$(vh_cmd) -- mkdir -p /filestorage/logs
#####################################################################################
copy_models: 
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/ac/graph_w63g /models/ac/
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/ac/19_11_18-3g /models/ac/

copy_rnnlm: 
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/rnnlm/v1 /models/rnnlm/

copy_lm: 
	./volume-helper/copyModel.sh $(LOCAL_MODELS_DIR)/lm/v1 /models/lm/
########### APPS ####################################################################
copy_apps: 
	./volume-helper/copyModel.sh $(KALDI_APPS_DIR) /apps/
#####################################################################################
redeploy:
	kubectl apply -f deployments/

clean:

.PHONY:
	ku_po init_vh stop_vh test_upload
