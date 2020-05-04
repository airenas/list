-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=metricsCollector
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=metrics-collector
tag=airenas/$(docker_container_prefix)$(service):$(METRICS_COLLECTOR_VERSION).$(commit_count)

#####################################################################################
$(dist_dir):
	mkdir -p $@
$(dist_dir)/$(executable_name): $(GO_SRC_DIR) | $(dist_dir)
	cd $(main_dir) && $(go_build_cmd) -o $(dist_dir)/$(executable_name)
	@echo -n "$(GREEN)Build:$(C_NC) $(executable_name)\n"
build: $(dist_dir)/$(executable_name)
########### DOCKER ##################################################################
dbuild: build $(dist_script_files)
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
