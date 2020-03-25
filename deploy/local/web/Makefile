-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
main_dir=../../../src/web/ear-app
service=ear-app
version=$(EAR_APP_VERSION)
#####################################################################################
$(dist_dir):
	mkdir -p $(dist_dir)

updateVersion: | $(dist_dir)
	cat $(main_dir)/src/environments/environment.prod.sample \
		| sed 's/BUILD_VERSION/$(version).$(commit_count)/g' > $(dist_dir)/environment.prod.ts
	rsync --checksum $(dist_dir)/environment.prod.ts $(main_dir)/src/environments/environment.prod.ts	

$(dist_dir)/.build: $(main_dir) $(main_dir)/src $(main_dir)/src/environments/environment.prod.ts
	cd $(main_dir) && ng build --prod --output-path=$(dist_dir)/html
	touch $(dist_dir)/.build

build: updateVersion $(dist_dir)/.build
########### DOCKER ##################################################################
tag=airenas/$(docker_container_prefix)$(service):$(version).$(commit_count)

dbuild: build
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
########### END #####################################################################
#####################################################################################
redeploy: build
	$(call dc-deploy,..,web)
#####################################################################################
clean:
	rm -rf $(dist_dir)

.PHONY:
	clean dbuild dpush build redeploy updateVersion
