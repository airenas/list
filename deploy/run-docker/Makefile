include Makefile.options
########################################################################################################
dwn_dir=$(deploy_dir)/dwn
install_dir=$(deploy_dir)/install
rec_file=$(deploy_dir)/volumes/models/config/recognizers.map.yml
dwn_url=https://prn509.vdu.lt:7080/testdata
apps_url=$(dwn_url)/apps-1.0.tar.gz
########################################################################################################
models_dwn=$(patsubst %, $(install_dir)/.%.done.dwn, $(models))
models_conf=$(patsubst %, $(install_dir)/.%.done.conf, $(models))
########################################################################################################
$(deploy_dir):
	mkdir -p $@
$(deploy_dir)/volumes:
	mkdir -p $@	
$(deploy_dir)/volumes/models:
	mkdir -p $@		
$(deploy_dir)/volumes/models/config:
	mkdir -p $@
$(dwn_dir):
	mkdir -p $@	
$(install_dir):
	mkdir -p $@		
########################################################################################################
$(install_dir)/.init.done: validate | $(deploy_dir)
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
	@$(if $(strip $(https_port)),echo "https_port = $(https_port)",echo No https_port set && exit 1)


########################################################################################################
$(install_dir)/.env.done: $(install_dir)/.init.done validate
	echo "LIST_TRAEFIK_HTTP_EXT_PORT=$(http_port)" >> $(deploy_dir)/.env 
	echo "LIST_TRAEFIK_HTTPS_EXT_PORT=$(https_port)" >> $(deploy_dir)/.env 
	echo "LIST_RABBITMQ_PASS=$(rabbitmq_pass)" >> $(deploy_dir)/.env 
	echo "LIST_MONGO_PASS=$(mongo_pass)" >> $(deploy_dir)/.env 
	echo "LIST_MONGO_URL=mongodb://list:$(mongo_pass)@mongo:27017" >> $(deploy_dir)/.env 
	echo "LIST_VOLUMES_DIR=$(deploy_dir)/volumes" >> $(deploy_dir)/.env 

	echo "HOST_EXTERNAL_URL=$(host_external_url)" >> $(deploy_dir)/.env 
	echo "SMTP_HOST=$(smtp_host)" >> $(deploy_dir)/.env 
	echo "SMTP_PORT=$(smtp_port)" >> $(deploy_dir)/.env 
	echo "SMTP_USERNAME=$(smtp_username)" >> $(deploy_dir)/.env 
	echo "SMTP_PASSWORD=$(smtp_password)" >> $(deploy_dir)/.env 
	echo "MAIL_URL=$(host_external_url)/ausis/results/{{ID}}" >> $(deploy_dir)/.env 

	touch $@

$(install_dir)/.dir.done: $(install_dir)/.init.done
	mkdir -p $(deploy_dir)/volumes/mongo
	mkdir -p $(deploy_dir)/volumes/rabbitmq
	mkdir -p $(deploy_dir)/volumes/fileStorage
	touch $@

########################################################################################################
$(dwn_dir)/apps.tar.gz: | $(dwn_dir)
	wget --no-check-certificate -O $(dwn_dir)/apps.tar.gz $(apps_url)

$(install_dir)/.apps.done.dwn: $(dwn_dir)/apps.tar.gz | $(deploy_dir)/volumes
	tar xvzf $(dwn_dir)/apps.tar.gz -C $(deploy_dir)/volumes
	touch $@
########################################################################################################
$(dwn_dir)/%.tar.gz: | $(dwn_dir)
	wget --no-check-certificate -O $@ $(dwn_url)/$*.tar.gz	

$(install_dir)/.%.done.dwn: $(dwn_dir)/%.tar.gz | $(install_dir) $(deploy_dir)/volumes/models
	tar xvzf $^ -C $(deploy_dir)/volumes/models
	touch $@
########################################################################################################
$(rec_file): $(deploy_dir)/volumes/models/config
	touch $@
$(install_dir)/.ben-tel.done.conf: | $(install_dir) $(rec_file)
	grep -qxF 'ben-tel: words_v3' $(rec_file) || echo 'ben-tel: words_v3' >> $(rec_file) 
	touch $@
$(install_dir)/.ben.done.conf: | $(install_dir) $(rec_file)
	grep -qxF 'ben: words_v3' $(rec_file) || echo 'ben: words_v3' >> $(rec_file) 
	touch $@
$(install_dir)/.adm-tel.done.conf: | $(install_dir) $(rec_file)
	grep -qxF 'adm-tel: words_v3' $(rec_file) || echo 'adm-tel: words_v3' >> $(rec_file) 
	touch $@
$(install_dir)/.adm.done.conf: | $(install_dir) $(rec_file)
	grep -qxF 'adm: words_v3' $(rec_file) || echo 'adm: words_v3' >> $(rec_file) 
	touch $@
$(install_dir)/.med-tel.done.conf: | $(install_dir) $(rec_file)
	grep -qxF 'med-tel: words_v3' $(rec_file) || echo 'med-tel: words_v3' >> $(rec_file) 
	touch $@
$(install_dir)/.med.done.conf: | $(install_dir) $(rec_file)
	grep -qxF 'med: words_v3' $(rec_file) || echo 'med: words_v3' >> $(rec_file) 
	touch $@
$(install_dir)/.tei-tel.done.conf: | $(install_dir) $(rec_file)
	grep -qxF 'tei-tel: words_v3' $(rec_file) || echo 'tei-tel: words_v3' >> $(rec_file) 
	touch $@
$(install_dir)/.tei.done.conf: | $(install_dir) $(rec_file)
	grep -qxF 'tei: words_v3' $(rec_file) || echo 'tei: words_v3' >> $(rec_file) 
	touch $@	
########################################################################################################
$(deploy_dir)/install/.ben-tel.done.dwn: $(install_dir)/.lm-1.0.done.dwn $(install_dir)/.ac-1.0.done.dwn 
	touch $@
$(deploy_dir)/install/.ben.done.dwn: $(install_dir)/.lm-1.0.done.dwn $(install_dir)/.ac-1.0.done.dwn 
	touch $@
$(deploy_dir)/install/.adm-tel.done.dwn: $(install_dir)/.lm-1.0.done.dwn $(install_dir)/.ac-1.0.done.dwn 
	touch $@
$(deploy_dir)/install/.adm.done.dwn: $(install_dir)/.lm-1.0.done.dwn $(install_dir)/.ac-1.0.done.dwn 
	touch $@
$(deploy_dir)/install/.tei-tel.done.dwn: $(install_dir)/.lm-1.0.done.dwn $(install_dir)/.ac-1.0.done.dwn 
	touch $@
$(deploy_dir)/install/.tei.done.dwn: $(install_dir)/.lm-1.0.done.dwn $(install_dir)/.ac-1.0.done.dwn 
	touch $@
$(deploy_dir)/install/.med-tel.done.dwn: $(install_dir)/.lm-1.0.done.dwn $(install_dir)/.ac-1.0.done.dwn 
	touch $@
$(deploy_dir)/install/.med.done.dwn: $(install_dir)/.lm-1.0.done.dwn $(install_dir)/.ac-1.0.done.dwn 
	touch $@		
########################################################################################################
init: $(install_dir)/.dir.done $(install_dir)/.env.done $(install_dir)/.prom.done $(install_dir)/.apps.done.dwn \
	$(install_dir)/.punctuation-1.0.done.dwn $(install_dir)/.diarization-1.0.done.dwn $(models_dwn) $(models_conf)

$(install_dir)/.prom.done: $(install_dir)/.env.done
	mkdir -p $(deploy_dir)/volumes/prometheus
	cp prometheus/* $(deploy_dir)/volumes/prometheus
	touch $@
########################################################################################################
run: init
	cd $(deploy_dir) && docker-compose up -d

install: run	
########################################################################################################
status: 
	cd $(deploy_dir) && docker-compose ps
########################################################################################################
clean: clean-docker clean-dir 

clean-dir:
	rm -rf $(deploy_dir)	

clean-docker:
	cd $(deploy_dir) && docker-compose stop && docker-compose rm -f
########################################################################################################
.SECONDARY:

.PHONY: validate clean run init install