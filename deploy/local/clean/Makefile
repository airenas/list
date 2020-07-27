-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=cleanService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=clean-service
version=$(CLEAN_SERVICE_VERSION)
#####################################################################################
$(dist_dir)/$(executable_name): $(GO_SRC_DIR)
	cd $(main_dir) && CGO_ENABLED=0 $(go_build_cmd) -o $@

build: $(dist_dir)/$(executable_name)
#####################################################################################
$(dist_dir):
	mkdir -p $@
########### DOCKER ##################################################################
tag=$(docker_prefix)$(service):$(version).$(commit_count)
dbuild: $(dist_dir)/$(executable_name)
	docker build -t $(tag) ./
dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
########### END #####################################################################
#####################################################################################
redeploy: build
	$(call dc-deploy,..,$(service))
#####################################################################################

clean:
	rm -rf $(dist_dir)

.PHONY:
	clean copy build dbuild dpush
