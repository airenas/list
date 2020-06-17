include Makefile.options
########################################################################################################
dwn_dir=$(deploy_dir)/dwn
install_dir=$(deploy_dir)/install
dwn_url=https://prn509.vdu.lt:7080/testdata
apps_url=$(dwn_url)/apps-1.0.tar.gz
punctuation_url=$(dwn_url)/punctuation-1.0.tar.gz
diarization_url=$(dwn_url)/diarization-1.0.tar.gz
########################################################################################################
models_dot=$(patsubst %, $(deploy_dir)/install/.%.model.done, $(models))
########################################################################################################
$(deploy_dir):
	mkdir -p $@

$(deploy_dir)/volumes:
	mkdir -p $@	

$(dwn_dir):
	mkdir -p $@	

$(install_dir):
	mkdir -p $@		
########################################################################################################
$(deploy_dir)/.init.done: | $(deploy_dir)
	cp docker-compose.yml $(deploy_dir)/docker-compose.yml
	cp .env $(deploy_dir)/.env
	touch $@

$(deploy_dir)/.env.done: $(deploy_dir)/.init.done
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

$(deploy_dir)/.dir.done: $(deploy_dir)/.init.done
	mkdir -p $(deploy_dir)/volumes/mongo
	mkdir -p $(deploy_dir)/volumes/rabbitmq
	mkdir -p $(deploy_dir)/volumes/fileStorage
	touch $@

########################################################################################################
$(dwn_dir)/apps.tar.gz: | $(dwn_dir)
	wget --no-check-certificate -O $(dwn_dir)/apps.tar.gz $(apps_url)

$(deploy_dir)/.apps.done: $(dwn_dir)/apps.tar.gz | $(deploy_dir)/volumes
	tar xvzf $(dwn_dir)/apps.tar.gz -C $(deploy_dir)/volumes
	touch $@
########################################################################################################
$(dwn_dir)/%.tar.gz: | $(dwn_dir)
	wget --no-check-certificate -O $@ $(dwn_url)/$*.tar.gz	

$(install_dir)/.%.done: $(dwn_dir)/%.tar.gz | $(install_dir)
	tar xvzf $^ -C $(deploy_dir)/volumes/models
	touch $@
########################################################################################################
$(deploy_dir)/install/.ben-tel.model.done: $(install_dir)/.lm-1.0.done $(install_dir)/.ac-1.0.done 
	touch $@
$(deploy_dir)/install/.ben.model.done: $(install_dir)/.lm-1.0.done $(install_dir)/.ac-1.0.done 
	touch $@
$(deploy_dir)/install/.adm-tel.model.done: $(install_dir)/.lm-1.0.done $(install_dir)/.ac-1.0.done 
	touch $@
$(deploy_dir)/install/.adm.model.done: $(install_dir)/.lm-1.0.done $(install_dir)/.ac-1.0.done 
	touch $@
$(deploy_dir)/install/.tei-tel.model.done: $(install_dir)/.lm-1.0.done $(install_dir)/.ac-1.0.done 
	touch $@
$(deploy_dir)/install/.tei.model.done: $(install_dir)/.lm-1.0.done $(install_dir)/.ac-1.0.done 
	touch $@
$(deploy_dir)/install/.med-tel.model.done: $(install_dir)/.lm-1.0.done $(install_dir)/.ac-1.0.done 
	touch $@
$(deploy_dir)/install/.med.model.done: $(install_dir)/.lm-1.0.done $(install_dir)/.ac-1.0.done 
	touch $@		
########################################################################################################
init: $(deploy_dir)/.dir.done $(deploy_dir)/.env.done $(deploy_dir)/.prom.done $(deploy_dir)/.apps.done \
	$(install_dir)/.punctuation-1.0.done $(install_dir)/.diarization-1.0.done $(models_dot)
	
$(deploy_dir)/.prom.done: $(deploy_dir)/.env.done
	mkdir -p $(deploy_dir)/volumes/prometheus
	cp prometheus/* $(deploy_dir)/volumes/prometheus
	touch $@
########################################################################################################

run: init
	cd $(deploy_dir) && docker-compose up -d

clean: clean-docker clean-dir 

clean-dir:
	rm -rf $(deploy_dir)	

clean-docker:
	cd $(deploy_dir) && docker-compose stop && docker-compose rm -f
