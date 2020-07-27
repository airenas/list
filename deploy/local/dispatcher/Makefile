-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=dispatcherService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=dispatcher-service
#####################################################################################
$(dist_dir):
	mkdir -p $@

$(dist_dir)/$(executable_name): $(GO_SRC_DIR) | $(dist_dir)
	cd $(main_dir) && $(go_build_cmd) -o $(dist_dir)/$(executable_name)

build: $(dist_dir)/$(executable_name) 
########### DOCKER ##################################################################
tag=$(docker_prefix)$(service):$(DISPATCHER_SERVICE_VERSION).$(commit_count)

dbuild: $(dist_dir)/$(executable_name)
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
########### END #####################################################################
#####################################################################################
redeploy: build 
	$(call dc-deploy,..,transcription-$(service))
#####################################################################################
clean:
	rm -rf $(dist_dir)

.PHONY:
	clean copy build dbuild dpush
