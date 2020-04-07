-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=cmdWorkerService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=result-make-service
tag=airenas/$(docker_container_prefix)$(service):$(TR_RESULT_SERVICE_VERSION).$(commit_count)

punct_wrapper_name=punct.file
punct_wrapper_dir=$(SRC_DIR)/tools/$(punct_wrapper_name)

pipe_runner_name=pipe.runner
pipe_runner_dir=$(SRC_DIR)/tools/$(pipe_runner_name)
#####################################################################################
$(dist_dir):
	mkdir -p $@
$(dist_dir)/$(executable_name): $(GO_SRC_DIR) | $(dist_dir)
	cd $(main_dir) && $(go_build_cmd) -o $(dist_dir)/$(executable_name)
	@echo -n "$(GREEN)Build:$(C_NC) $(executable_name)\n"
$(dist_dir)/$(punct_wrapper_name): | $(dist_dir)
	cd $(punct_wrapper_dir) && $(go_build_cmd) -o $(dist_dir)/$(punct_wrapper_name)
	@echo -n "$(GREEN)Build:$(C_NC) $(punct_wrapper_name)\n"
$(dist_dir)/$(pipe_runner_name): | $(dist_dir)
	cd $(pipe_runner_dir) && $(go_build_cmd) -o $(dist_dir)/$(pipe_runner_name)
	@echo -n "$(GREEN)Build:$(C_NC) $(pipe_runner_name)\n"
build: $(dist_dir)/$(executable_name) $(dist_dir)/$(pipe_runner_name) $(dist_dir)/$(punct_wrapper_name)
#####################################################################################
$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/result.make/Makefile
	mkdir -p $(@D)
	cp $< $@
copy: $(dist_dir)/cmd/Makefile
########### DOCKER ##################################################################
dbuild: build copy
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
########### END #####################################################################
#####################################################################################
redeploy: build copy
	$(call dc-deploy,..,$(service))
#####################################################################################

clean:
	rm -rf $(dist_dir)

.PHONY:
	clean copy build dbuild dpush
