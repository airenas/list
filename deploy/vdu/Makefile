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

remote_dwn_dir=$(remote_volumes_dir)/dwn
remote_models_dir=$(remote_volumes_dir)/models
dwn_url?=https://prn509.vdu.lt:7080/testdata
#####################################################################################
voice_files=$(dist_dir)/voice/lab_ALL.htsvoice
$(dist_dir)/voice/%: $(voice_dir)/%  
	mkdir -p $(@D)
	cp $< $@
########### CONFIG ##################################################################
copy-config: 
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
	@$(if $(strip $(srv)),echo "Service = $(srv)",echo No srv set && exit 1)
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose up --detach --build $(srv)"

log-one:
	@$(if $(strip $(srv)),echo "Service = $(srv)",echo No srv set && exit 1)
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose logs $(srv)"

exec-into:
	@$(if $(strip $(srv)),echo "Service = $(srv)",echo No srv set && exit 1)
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose exec $(srv) /bin/bash"

status:
	$(SSH_CMD) "cd $(HOME_DIR)/$(remote_service_dir) && docker-compose ps"
ssh:
	$(SSH_CMD)

tunel-mongo:
	ssh -L 27018:localhost:27017 -t -p $(SSH_PORT) $(USER)@$(SERVER_URL)	
########### END SERVICES ############################################################
########### MODELS ##################################################################
install-models: install-am-ben-2.0.data install-config-ben-2.0.data
install-rnnlm: install-rnnlm-2.0.data
install-lm: install-lm-ben-2.0.data
copy_diarization: 	
	$(RSYNC_S_CMD) $(LOCAL_MODELS_DIR)/diarization/v1/* $(USER)@$(SERVER_URL):$(remote_volumes_dir)/models/diarization/	
install-punctuation: install-punctuation-2.0.data
########### DWN #####################################################################
fake-%.dwn: 
	$(SSH_CMD) "mkdir -p $(remote_dwn_dir)"
	$(SSH_CMD) "wget --no-check-certificate -O $(remote_dwn_dir)/$*.tar.gz $(dwn_url)/$*.tar.gz"
install-%.data: fake-%.dwn
	$(SSH_CMD) "tar xvzf $(remote_dwn_dir)/$*.tar.gz -C $(remote_models_dir)"
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
