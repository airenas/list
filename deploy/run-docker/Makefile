include Makefile.options

$(deploy_dir):
	mkdir -p $@

$(deploy_dir)/volumes:
	mkdir -p $@	

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
	mkdir -p $(deploy_dir)/volumes/
	touch $@

init: $(deploy_dir)/.dir.done $(deploy_dir)/.env.done

run: $(deploy_dir)/.env.done $(deploy_dir)/.dir.done
	cd $(deploy_dir) && docker-compose up -d

clean: clean-docker clean-dir 

clean-dir:
	rm -rf $(deploy_dir)	

clean-docker:
	cd $(deploy_dir) && docker-compose stop && docker-compose rm -f
