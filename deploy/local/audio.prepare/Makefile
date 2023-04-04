-include ../Makefile.options
-include ../../../version
#####################################################################################
dist_dir=$(CURDIR)/deploy
service=audio-prepare-service
denoise_dir=../../../sfgparse
#####################################################################################
tools_dir=$(SRC_DIR)/tools
tools=m_shell.sh m_start.sh m_end.sh sfg
dist_tools_files=$(patsubst %, $(dist_dir)/%, $(tools))
#####################################################################################
scripts_dir=$(SRC_DIR)/decode/scripts
#####################################################################################
build: $(dist_tools_files)
#####################################################################################
$(dist_dir)/cmd:
	mkdir -p $(dist_dir)/cmd
$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/audio.prepare/Makefile | $(dist_dir)/cmd
	cp $< $@
$(dist_dir)/%.sh: $(scripts_dir)/%.sh | $(dist_dir)/cmd
	cp $< $@
$(dist_dir)/sfg: $(denoise_dir)/sfg | $(dist_dir)
	cp $< $@
copy: $(dist_dir)/cmd/Makefile $(dist_dir)/$(executable_name) $(dist_tools_files) $(dist_scripts_files)
########### DOCKER ##################################################################
tag=$(docker_prefix)$(service):$(audio_prepare_version)

dbuild: $(dist_dir)/cmd/Makefile $(dist_tools_files) $(dist_scripts_files)
	cd ../../../ && docker buildx build -t $(tag) --build-arg BUILD_VERSION=$(version) --build-arg TOOLS_VERSION=$(tools_version) \
		-f deploy/local/audio.prepare/Dockerfile ./

dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
########### END #####################################################################
clean:
	rm -rf $(dist_dir)

.PHONY:
	clean copy build dbuild dpush
