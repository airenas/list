-include ../Makefile.options
#####################################################################################
dist_dir=deploy
executable_name=statusProviderService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=status-service
version=$(STATUS_SERVICE_VERSION)
#####################################################################################
$(main_dir)/$(executable_name): $(GO_SRC_DIR)
	(cd $(main_dir) && $(go_build_cmd))

build: $(main_dir)/$(executable_name)
#####################################################################################
$(dist_dir)/$(executable_name): $(main_dir)/$(executable_name)
	mkdir -p $(@D)
	cp $< $@

copy: $(dist_dir)/$(executable_name)
########### DOCKER ##################################################################
tag=airenas/$(docker_container_prefix)$(service):$(version).$(commit_count)

dbuild: $(dist_dir)/$(executable_name)
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
########### END #####################################################################
#####################################################################################
redeploy: build copy
	$(call dc-deploy,..,$(service))
#####################################################################################
clean:
	rm -rf $(dist_dir)
	rm -f $(main_dir)/$(executable_name)

.PHONY:
	clean copy build dbuild dpush redeploy
