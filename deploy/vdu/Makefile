-include Makefile.options
#####################################################################################
USER?=airenas
SERVER_URL?=list.airenas.eu
SSH_PORT?=7022
HOME_DIR?=/home/airenas

SSH_CMD=ssh -p $(SSH_PORT) $(USER)@$(SERVER_URL)
RSYNC_S_CMD=rsync -r --size-only -e 'ssh -p $(SSH_PORT)' -P 
RSYNC_CMD=rsync -r -e 'ssh -p $(SSH_PORT)' -P 
#####################################################################################
src_dir=../../
#####################################################################################
remote_service_dir=services
remote_volumes_dir=$(HOME_DIR)/list.volumes
remote_acme_file=$(HOME_DIR)/$(remote_service_dir)/traefik/acme.json
#####################################################################################
voice_files=$(dist_dir)/voice/lab_ALL.htsvoice
$(dist_dir)/voice/%: $(voice_dir)/%  
	mkdir -p $(@D)
	cp $< $@
#####################################################################################
copy_local: 
	rsync -r -e -P $(src_dir)/decode/models/words_rnnlm/scripts models/config/words_v2

########### CONFIG ##################################################################
copy_config: 
	$(RSYNC_CMD) models $(USER)@$(SERVER_URL):$(remote_volumes_dir)/

########### SERVICES ################################################################
copy: 
	$(RSYNC_CMD) $(remote_service_dir) $(USER)@$(SERVER_URL):$(HOME_DIR)/
	$(SSH_CMD) "cp -n $(remote_acme_file).sample $(remote_acme_file)"
	$(SSH_CMD) "chmod 600 $(remote_acme_file)"

deploy: copy
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose stop && docker-compose build && docker-compose pull && docker-compose up -d"

status:
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose ps"
########### END SERVICES ############################################################
########### MODELS ##################################################################
copy_models: 
	$(RSYNC_S_CMD) -P $(LOCAL_MODELS_DIR)/ac/phones_v0 $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/ac/
	$(RSYNC_S_CMD) -P $(LOCAL_MODELS_DIR)/ac/graph_w63g $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/ac/
	$(RSYNC_S_CMD) -P $(LOCAL_MODELS_DIR)/ac/19_11_18-3g $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/ac/

########### END #####################################################################
clean:
	rm -rf $(dist_dir)

.PHONY:
	copy deploy clean copy_local status copy_config
