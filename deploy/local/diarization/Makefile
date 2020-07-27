-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=cmdWorkerService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=diarization-service
version=$(DIARIZATION_SERVICE_VERSION)
#####################################################################################
tools=send.metric fix.segments \
	m_shell.sh m_start.sh m_end.sh
tools_dir=$(SRC_DIR)/tools
dist_tools_files=$(patsubst %, $(dist_dir)/%, $(tools))
#####################################################################################
$(dist_dir):
	mkdir -p $@
$(dist_dir)/$(executable_name): $(GO_SRC_DIR) | $(dist_dir)
	cd $(main_dir) && CGO_ENABLED=0 $(go_build_cmd) -o $(dist_dir)/$(executable_name)
$(dist_dir)/%: $(tools_dir)/% | $(dist_dir)
	cd $< && $(go_build_cmd) -o $@
	@echo -n "$(GREEN)Build:$(C_NC) $*\n"
build: $(dist_dir)/$(executable_name) $(dist_tools_files)
#####################################################################################
$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/diarization/Makefile
	mkdir -p $(@D)
	cp $< $@
$(dist_dir)/%.sh: $(SRC_DIR)/decode/scripts/%.sh | $(dist_dir)
	cp $< $@
copy: $(dist_dir)/cmd/Makefile $(dist_dir)/$(executable_name)
########### DOCKER ##################################################################
tag=$(docker_prefix)$(service):$(version).$(commit_count)

dbuild: $(dist_dir)/cmd/Makefile $(dist_dir)/$(executable_name) $(dist_tools_files)
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
########### END #####################################################################
clean:
	rm -rf $(dist_dir)
	rm -f $(main_dir)/$(executable_name)

.PHONY:
	clean copy build dbuild dpush
