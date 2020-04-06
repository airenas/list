-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=cmdWorkerService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=transcription-service

pipe_runner_name=pipe.runner
pipe_runner_dir=$(SRC_DIR)/tools/$(pipe_runner_name)
#####################################################################################
$(dist_dir):
	mkdir -p $@

$(dist_dir)/$(executable_name): $(GO_SRC_DIR) | $(dist_dir)
	cd $(main_dir) && $(go_build_cmd) -o $(dist_dir)/$(executable_name)

$(dist_dir)/$(pipe_runner_name): | $(dist_dir)
	cd $(pipe_runner_dir) && $(go_build_cmd) -o $(dist_dir)/$(pipe_runner_name)

build: $(dist_dir)/$(executable_name) $(dist_dir)/$(pipe_runner_name)
#####################################################################################
$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/transcription/Makefile | $(dist_dir)
	mkdir -p $(@D)
	cp $< $@

copy: $(dist_dir)/cmd/Makefile 
########### DOCKER ##################################################################
tag=airenas/$(docker_container_prefix)$(service):$(TRANSCRIPTION_SERVICE_VERSION).$(commit_count)

dbuild: $(dist_dir)/cmd/Makefile $(dist_dir)/$(executable_name) $(dist_dir)/$(pipe_runner_name)
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
	rm -f $(main_dir)/$(executable_name)

.PHONY:
	clean copy build dbuild dpush
