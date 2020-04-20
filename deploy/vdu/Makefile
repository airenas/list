-include Makefile.options
#####################################################################################
USER?=airenas
SERVER_URL?=list.airenas.eu
SSH_PORT?=7022
HOME_DIR?=/home/airenas

SSH_CMD=ssh -t -p $(SSH_PORT) $(USER)@$(SERVER_URL)
RSYNC_S_CMD=rsync -r --size-only -e 'ssh -p $(SSH_PORT)' -P 
RSYNC_CMD=rsync -r -e 'ssh -p $(SSH_PORT)' -P 
#####################################################################################
src_dir=../../src
apps_dir=../local/kaldi/apps
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

init_volumes: 
	$(SSH_CMD) "mkdir -p $(remote_volumes_dir)/mongo"
	$(SSH_CMD) "mkdir -p $(remote_volumes_dir)/rabbitmq"
	$(SSH_CMD) "mkdir -p $(remote_volumes_dir)/fileStorage"
	$(SSH_CMD) "mkdir -p $(remote_volumes_dir)/fileStorage/decoded/audio"
	$(SSH_CMD) "mkdir -p $(remote_volumes_dir)/fileStorage/decoded/trans"
	$(SSH_CMD) "mkdir -p $(remote_volumes_dir)/fileStorage/results"
	$(SSH_CMD) "mkdir -p $(remote_volumes_dir)/fileStorage/logs"

deploy: copy
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose stop && docker-compose build && docker-compose pull && docker-compose up -d"

deploy-one: copy
	@$(if $(strip $(service)),echo "Service = $(service)",echo No service set && exit 1)
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose stop $(service) && docker-compose pull $(service) && docker-compose up -d $(service)"

log-one:
	@$(if $(strip $(service)),echo "Service = $(service)",echo No service set && exit 1)
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose logs $(service)"

exec-into:
	@$(if $(strip $(service)),echo "Service = $(service)",echo No service set && exit 1)
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose exec $(service) /bin/bash"


status:
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose ps"
ssh:
	$(SSH_CMD)
########### END SERVICES ############################################################
########### MODELS ##################################################################
copy_models: 
	$(RSYNC_S_CMD) $(LOCAL_MODELS_DIR)/ac/phones_v0 $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/ac/
	$(RSYNC_S_CMD) $(LOCAL_MODELS_DIR)/ac/graph_w63g $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/ac/
	$(RSYNC_S_CMD) $(LOCAL_MODELS_DIR)/ac/19_11_18-3g $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/ac/
copy_rnnlm: 
	$(RSYNC_S_CMD) $(LOCAL_MODELS_DIR)/rnnlm/v1 $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/rnnlm/
copy_lm: 
	$(RSYNC_S_CMD) $(LOCAL_MODELS_DIR)/lm/v1 $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/lm/
copy_diarization: 
	$(RSYNC_S_CMD) $(LOCAL_MODELS_DIR)/diarization/v1/* $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/diarization/	
copy_punctuation: 
	$(RSYNC_S_CMD) $(LOCAL_MODELS_DIR)/punctuation/h256_w200k/* $(USER)@$(SERVER_URL):$(remote_volumes_dir)/h256_w200k/		
########### END #####################################################################
########### APPS ####################################################################
copy_apps: 
	$(RSYNC_S_CMD) $(apps_dir) $(USER)@$(SERVER_URL):$(remote_volumes_dir)/
########### END #####################################################################
copy_tmp: 
	$(RSYNC_S_CMD) ../../../listgo/example/cmd.run/cmd.run $(USER)@$(SERVER_URL):$(remote_volumes_dir)/apps/kaldi/bin

clean:
	rm -rf $(dist_dir)

.PHONY:
	copy deploy clean copy_local status copy_config
