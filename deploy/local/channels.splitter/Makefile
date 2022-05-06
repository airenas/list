-include ../Makefile.options
-include ../../../version
#####################################################################################
dist_dir=$(CURDIR)/deploy
service=channels-splitter
version=$(channels_splitter_version)
#####################################################################################
$(dist_dir):
	mkdir -p $@
#####################################################################################
tools=m_shell.sh m_start.sh m_end.sh
dist_tools_files=$(patsubst %, $(dist_dir)/%, $(tools))
#####################################################################################
scripts_dir=$(SRC_DIR)/decode/scripts
#####################################################################################
$(dist_dir)/cmd:
	mkdir -p $(dist_dir)/cmd
$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/channels.split/Makefile | $(dist_dir)/cmd
	cp $< $@
$(dist_dir)/%.sh: $(scripts_dir)/%.sh | $(dist_dir)/cmd
	cp $< $@
copy: $(dist_dir)/cmd/Makefile $(dist_tools_files) $(dist_scripts_files)
########### DOCKER ##################################################################
tag=$(docker_prefix)$(service):$(version)

dbuild: $(dist_dir)/cmd/Makefile $(dist_tools_files) $(dist_scripts_files)
	cd ../../../ && docker build -t $(tag) --build-arg BUILD_VERSION=$(version) --build-arg TOOLS_VERSION=$(tools_version) \
		-f deploy/local/channels.splitter/Dockerfile ./

dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
dscan: dbuild
	docker scan --accept-license $(tag)			
########### END #####################################################################
clean:
	rm -rf $(dist_dir)

.PHONY:
	clean copy build dbuild dpush dscan
