-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=punctuationService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=punctuation-service
version=$(PUNCTUATION_SERVICE_VERSION)
#####################################################################################
$(dist_dir):
	mkdir -p $@
$(dist_dir)/$(executable_name): $(GO_SRC_DIR)
	cd $(main_dir) && $(go_build_cmd) -o $@

build: $(dist_dir)/$(executable_name)
########### DOCKER ##################################################################
tag=$(docker_prefix)$(service):$(version).$(commit_count)

dbuild: $(dist_dir)/$(executable_name)
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
########### END #####################################################################
#####################################################################################
redeploy: $(dist_dir)/$(executable_name)
	$(call dc-deploy,..,$(service))
#####################################################################################
clean:
	rm -rf $(dist_dir)
	rm -f $(main_dir)/$(executable_name)

.PHONY:
	clean copy build dbuild dpush redeploy

