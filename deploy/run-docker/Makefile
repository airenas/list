include Makefile.options
########################################################################################################
dwn_dir=$(deploy_dir)/dwn
install_dir=$(deploy_dir)/.install
rec_file=$(deploy_dir)/volumes/models/config/recognizers.map.yml
dwn_url=https://semantikadocker.vdu.lt/models
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
$(install_dir)/.init.done: validate | $(deploy_dir) $(install_dir)
	cp docker-compose.yml $(deploy_dir)/docker-compose.yml
	cp .env $(deploy_dir)/.env
	touch $@
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
$(install_dir)/.env.done: $(install_dir)/.init.done validate
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

$(install_dir)/.%.data.extracted: $(install_dir)/.%.dwn | $(deploy_dir)/volumes/models
	tar xvzf $(dwn_dir)/$*.tar.gz -C $(deploy_dir)/volumes/models
	touch $@
########################################################################################################
$(rec_file): $(deploy_dir)/volumes/models/config
	touch $@
$(install_dir)/.ben-tel.conf: | $(install_dir) $(rec_file)
	grep -qxF 'ben-tel: ben-2.0' $(rec_file) || echo 'ben-tel: ben-2.0' >> $(rec_file) 
	touch $@
$(install_dir)/.ben.conf: | $(install_dir) $(rec_file)
	grep -qxF 'ben: ben-2.0' $(rec_file) || echo 'ben: ben-2.0' >> $(rec_file) 
	touch $@
$(install_dir)/.adm-tel.conf: | $(install_dir) $(rec_file)
	grep -qxF 'adm-tel: adm-1.0' $(rec_file) || echo 'adm-tel: adm-1.0' >> $(rec_file) 
	touch $@
$(install_dir)/.adm.conf: | $(install_dir) $(rec_file)
	grep -qxF 'adm: adm-1.0' $(rec_file) || echo 'adm: adm-1.0' >> $(rec_file) 
	touch $@
$(install_dir)/.med-tel.conf: | $(install_dir) $(rec_file)
	grep -qxF 'med-tel: med-1.0' $(rec_file) || echo 'med-tel: med-1.0' >> $(rec_file) 
	touch $@
$(install_dir)/.med.conf: | $(install_dir) $(rec_file)
	grep -qxF 'med: med-1.0' $(rec_file) || echo 'med: med-1.0' >> $(rec_file) 
	touch $@
$(install_dir)/.tei-tel.conf: | $(install_dir) $(rec_file)
	grep -qxF 'tei-tel: tei-1.0' $(rec_file) || echo 'tei-tel: tei-1.0' >> $(rec_file) 
	touch $@
$(install_dir)/.tei.conf: | $(install_dir) $(rec_file)
	grep -qxF 'tei: tei-1.0' $(rec_file) || echo 'tei: tei-1.0' >> $(rec_file) 
	touch $@	
########################################################################################################
$(install_dir)/.ben-tel.done: $(install_dir)/.lm-ben-2.0.data.extracted $(install_dir)/.am-ben-2.0.data.extracted $(install_dir)/.config-ben-2.0.data.extracted 
	touch $@
$(install_dir)/.ben.done: $(install_dir)/.lm-ben-2.0.data.extracted $(install_dir)/.am-ben-2.0.data.extracted $(install_dir)/.config-ben-2.0.data.extracted 
	touch $@
$(install_dir)/.adm-tel.done: $(install_dir)/.lm-adm-1.0.data.extracted $(install_dir)/.am-adm-1.0.data.extracted $(install_dir)/.config-adm-1.0.data.extracted 
	touch $@
$(install_dir)/.adm.done: $(install_dir)/.lm-adm-1.0.data.extracted $(install_dir)/.am-adm-1.0.data.extracted $(install_dir)/.config-adm-1.0.data.extracted 
	touch $@
$(install_dir)/.tei-tel.done: $(install_dir)/.lm-tei-1.0.data.extracted $(install_dir)/.am-tei-1.0.data.extracted $(install_dir)/.config-tei-1.0.data.extracted 
	touch $@
$(install_dir)/.tei.done: $(install_dir)/.lm-tei-1.0.data.extracted $(install_dir)/.am-tei-1.0.data.extracted $(install_dir)/.config-tei-1.0.data.extracted 
	touch $@
$(install_dir)/.med-tel.done: $(install_dir)/.lm-med-1.0.data.extracted $(install_dir)/.am-med-1.0.data.extracted $(install_dir)/.config-med-1.0.data.extracted 
	touch $@
$(install_dir)/.med.done: $(install_dir)/.lm-med-1.0.data.extracted $(install_dir)/.am-med-1.0.data.extracted $(install_dir)/.config-med-1.0.data.extracted 
	touch $@		
########################################################################################################
init: $(install_dir)/.dir.done $(install_dir)/.env.done $(install_dir)/.prom.done $(install_dir)/.apps-1.0.apps.extracted \
	$(install_dir)/.punctuation-2.0.data.extracted $(install_dir)/.diarization-1.0.data.extracted $(models_done) $(models_conf) \
	$(install_dir)/.rnnlm-2.0.data.extracted

$(install_dir)/.prom.done: $(install_dir)/.env.done
	mkdir -p $(deploy_dir)/volumes/prometheus
	cp prometheus/* $(deploy_dir)/volumes/prometheus
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