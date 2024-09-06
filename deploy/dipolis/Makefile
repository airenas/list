include Makefile.options
########################################################################################################
dwn_dir=$(deploy_dir)/dwn
install_dir=$(deploy_dir)/.install
rec_file=$(deploy_dir)/volumes/models/config/recognizers.map.yml
dwn_url=https://semantikadocker.vdu.lt/models
dwn2_url=https://sinteze.intelektika.lt/models/share-dipolis
########################################################################################################
models_done=$(patsubst %, $(install_dir)/.%.done, $(models))
models_conf=$(patsubst %, $(install_dir)/.%.conf, $(models))
########################################################################################################
$(deploy_dir):
	mkdir -p $@
$(deploy_dir)/volumes:
	mkdir -p $@	
$(deploy_dir)/volumes/models:
	mkdir -p $@		
$(deploy_dir)/volumes/models/config:
	mkdir -p $@
$(deploy_dir)/volumes/apps:
	mkdir -p $@	
$(dwn_dir):
	mkdir -p $@	
$(install_dir):
	mkdir -p $@		
########################################################################################################
$(install_dir)/.init.done: validate $(deploy_dir)/.env | $(deploy_dir) $(install_dir)
	cp docker-compose.yml $(deploy_dir)/docker-compose.yml
	touch $@
########################################################################################################
$(deploy_dir)/.env: .env.in | $(deploy_dir)
	cat $< | envsubst > $@
########################################################################################################
validate:
	@$(if $(strip $(deploy_dir)),echo "deploy_dir = $(deploy_dir)",echo No deploy_dir set && exit 1)
	@$(if $(strip $(models)),echo "models = $(models)",echo No models set && exit 1)
	@$(if $(strip $(rabbitmq_pass)),echo "rabbitmq_pass = ****",echo No rabbitmq_pass set && exit 1)
	@$(if $(strip $(mongo_pass)),echo "mongo_pass = ****",echo No mongo_pass set && exit 1)
	@$(if $(strip $(http_port)),echo "http_port = $(http_port)",echo No http_port set && exit 1)
	@$(if $(strip $(smtp_username)),echo "smtp_username = $(smtp_username)",echo No smtp_username set)
ifneq ($(strip $(smtp_username)),)
	@$(if $(strip $(smtp_password)),echo "smtp_password = ****",echo No smtp_password set && exit 1)
	@$(if $(strip $(host_external_url)),echo "host_external_url = $(host_external_url)",echo No host_external_url set && exit 1)
endif

########################################################################################################
$(install_dir)/.env.done: $(deploy_dir)/.env
	printf "\n" >> $(deploy_dir)/.env
	echo "LIST_TRAEFIK_HTTP_EXT_PORT=$(http_port)" >> $(deploy_dir)/.env 
	echo "LIST_RABBITMQ_PASS=$(rabbitmq_pass)" >> $(deploy_dir)/.env 
	echo "LIST_MONGO_PASS=$(mongo_pass)" >> $(deploy_dir)/.env 
	echo "LIST_MONGO_URL=mongodb://list:$(mongo_pass)@mongo:27017" >> $(deploy_dir)/.env 
	echo "LIST_VOLUMES_DIR=$(deploy_dir)/volumes" >> $(deploy_dir)/.env 

	echo "HOST_EXTERNAL_URL=$(host_external_url)" >> $(deploy_dir)/.env 
	echo "SMTP_HOST=$(smtp_host)" >> $(deploy_dir)/.env 
	echo "SMTP_PORT=$(smtp_port)" >> $(deploy_dir)/.env 
	echo "SMTP_USERNAME=$(smtp_username)" >> $(deploy_dir)/.env 
	echo "SMTP_PASSWORD=$(smtp_password)" >> $(deploy_dir)/.env 
	echo "SMTP_TYPE=$(smtp_type)" >> $(deploy_dir)/.env
	echo "MAIL_URL=$(host_external_url)/ausis/results/{{ID}}" >> $(deploy_dir)/.env 
ifeq ($(strip $(smtp_username)),)
	echo "SENDINFORMMESSAGES=false" >> $(deploy_dir)/.env 
else
	echo "SENDINFORMMESSAGES=true" >> $(deploy_dir)/.env 
endif
	echo "HF_API_TOKEN=$(hf_api_token)" >> $(deploy_dir)/.env
	touch $@

$(install_dir)/.dir.done: $(install_dir)/.init.done
	mkdir -p $(deploy_dir)/volumes/mongo
	mkdir -p $(deploy_dir)/volumes/rabbitmq
	mkdir -p $(deploy_dir)/volumes/fileStorage
	mkdir -p $(deploy_dir)/volumes/fileStorage/logs
	mkdir -p $(deploy_dir)/volumes/fileStorage/results
	touch $@

########################################################################################################
$(install_dir)/.%.apps.extracted: $(install_dir)/.%.dwn | $(deploy_dir)/volumes/apps
	tar xvzf $(dwn_dir)/$*.tar.gz -C $(deploy_dir)/volumes/apps
	touch $@
########################################################################################################
$(install_dir)/.%.dwn: | $(dwn_dir)
	wget --no-check-certificate -O $(dwn_dir)/$*.tar.gz $(dwn_url)/$*.tar.gz	
	touch $@

$(install_dir)/.%.dwn2: | $(dwn_dir)
	wget --user user --password $(share_pass) -O $(dwn_dir)/$*.tar.gz $(dwn2_url)/$*.tar.gz 
	touch $@

$(install_dir)/.%.data.extracted: $(install_dir)/.%.dwn | $(deploy_dir)/volumes/models
	tar xvzf $(dwn_dir)/$*.tar.gz -C $(deploy_dir)/volumes/models
	touch $@
$(install_dir)/.%.data2.extracted: $(install_dir)/.%.dwn2 | $(deploy_dir)/volumes/models
	tar xvzf $(dwn_dir)/$*.tar.gz -C $(deploy_dir)/volumes/models
	touch $@	
########################################################################################################
$(rec_file): $(deploy_dir)/volumes/models/config
	touch $@
$(install_dir)/.ben-tel-2.0r1.conf: | $(install_dir) $(rec_file)
	grep -qxF 'default: ben-tel-2.0r1' $(rec_file) || echo 'default: ben-tel-2.0r1' >> $(rec_file) 
	grep -qxF 'ben-tel: ben-tel-2.0r1' $(rec_file) || echo 'ben-tel: ben-tel-2.0r1' >> $(rec_file) 
	touch $@
########################################################################################################
$(install_dir)/.ben-tel-2.0r1.done: $(install_dir)/.lm-ben-tel-2.0.data2.extracted $(install_dir)/.am-ben-tel-2.0.data2.extracted $(install_dir)/.config-ben-tel-2.0r1.data2.extracted \
	$(install_dir)/.rnnlm-4.1.data2.extracted
	touch $@
########################################################################################################
init: $(install_dir)/.dir.done $(install_dir)/.env.done $(install_dir)/.prom.done $(install_dir)/.traefik.done $(install_dir)/.apps-1.0.apps.extracted \
	$(install_dir)/.punctuation-2.0.data.extracted $(install_dir)/.diarization-1.0.data.extracted $(models_done) $(models_conf) 

$(install_dir)/.prom.done: $(install_dir)/.env.done
	mkdir -p $(deploy_dir)/volumes/prometheus
	cp prometheus/* $(deploy_dir)/volumes/prometheus
	touch $@
$(install_dir)/.traefik.done: $(install_dir)/.env.done
	mkdir -p $(deploy_dir)/traefik
	cp -r traefik/* $(deploy_dir)/traefik
	touch $@	
########################################################################################################
run: init
	cd $(deploy_dir) && docker compose up -d

install: run	
########################################################################################################
status: 
	cd $(deploy_dir) && docker compose ps
status-service: 
	wget -SO- http://localhost:$(http_port)/ausis/transcriber/live?full=1 2>&1	
########################################################################################################
clean: clean-docker clean-dir 

clean-dir:
	rm -rf $(deploy_dir)	

clean-docker:
	cd $(deploy_dir) && docker compose down --rmi all
########################################################################################################
.PHONY: validate clean run init install
.EXPORT_ALL_VARIABLES: