-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=cmdWorkerService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=transcription-service

#####################################################################################
tools=send.metric pipe.runner \
	m_shell.sh m_start.sh m_end.sh
tools_dir=$(SRC_DIR)/tools
dist_tools_files=$(patsubst %, $(dist_dir)/%, $(tools))
#####################################################################################
$(dist_dir):
	mkdir -p $@

$(dist_dir)/$(executable_name): $(GO_SRC_DIR) | $(dist_dir)
	cd $(main_dir) && $(go_build_cmd) -o $(dist_dir)/$(executable_name)
$(dist_dir)/%: $(tools_dir)/% | $(dist_dir)
	cd $< && $(go_build_cmd) -o $@
	@echo -n "$(GREEN)Build:$(C_NC) $*\n"
build: $(dist_dir)/$(executable_name) $(dist_tools_files)
#####################################################################################
$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/transcription/Makefile | $(dist_dir)
	mkdir -p $(@D)
	cp $< $@
$(dist_dir)/%.sh: $(SRC_DIR)/decode/scripts/%.sh | $(dist_dir)
	cp $< $@

copy: $(dist_dir)/cmd/Makefile 
########### DOCKER ##################################################################
tag=$(docker_prefix)$(service):$(TRANSCRIPTION_SERVICE_VERSION).$(commit_count)

dbuild: $(dist_dir)/cmd/Makefile $(dist_dir)/$(executable_name) $(dist_tools_files)
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
