-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=cmdWorkerService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=audio-prepare-service
version=$(AUDIO_PREPARE_SERVICE_VERSION)
#####################################################################################
$(dist_dir):
	mkdir -p $@
$(dist_dir)/$(executable_name): $(GO_SRC_DIR) | $(dist_dir)
	cd $(main_dir) && CGO_ENABLED=0 $(go_build_cmd) -o $(dist_dir)/$(executable_name) 
	@echo -n "$(GREEN)Build:$(C_NC) $(executable_name)\n"
#####################################################################################
tools_dir=$(SRC_DIR)/tools
tools=send.metric m_shell.sh m_start.sh m_end.sh
dist_tools_files=$(patsubst %, $(dist_dir)/%, $(tools))
#####################################################################################
scripts_dir=$(SRC_DIR)/decode/scripts
#####################################################################################
$(dist_dir)/%: $(tools_dir)/% | $(dist_dir)
	cd $< && $(go_build_cmd) -o $@
	@echo -n "$(GREEN)Build:$(C_NC) $*\n"
build: $(dist_dir)/$(executable_name) $(dist_tools_files)
#####################################################################################
$(dist_dir)/cmd:
	mkdir -p $(dist_dir)/cmd
$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/audio.prepare/Makefile | $(dist_dir)/cmd
	cp $< $@
$(dist_dir)/%.sh: $(scripts_dir)/%.sh | $(dist_dir)/cmd
	cp $< $@
copy: $(dist_dir)/cmd/Makefile
########### DOCKER ##################################################################
tag=$(docker_prefix)$(service):$(version).$(commit_count)

dbuild: $(dist_dir)/cmd/Makefile $(dist_dir)/$(executable_name) $(dist_tools_files) $(dist_scripts_files)
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
