-include ../Makefile.options
-include ../../../version
#####################################################################################
dist_dir=$(CURDIR)/deploy
service=join-results-service
tag=$(docker_prefix)$(service):$(join_results_version)
#####################################################################################
tools=m_shell.sh m_start.sh m_end.sh 
dist_tools_files=$(patsubst %, $(dist_dir)/%, $(tools))
#####################################################################################
script_files=Makefile
dist_script_files=$(patsubst %, $(dist_dir)/cmd/%, $(script_files))
#####################################################################################
$(dist_dir):
	mkdir -p $@
#####################################################################################
$(dist_dir)/cmd: 
	mkdir -p $@
$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/join.results/Makefile | $(dist_dir)/cmd
	cp $< $@
$(dist_dir)/%.sh: $(SRC_DIR)/decode/scripts/%.sh | $(dist_dir)
	cp $< $@
copy: $(dist_script_files)
########### DOCKER ##################################################################
dbuild: $(dist_script_files) $(dist_tools_files)
	cd ../../../ && docker build -t $(tag) --build-arg BUILD_VERSION=$(version) --build-arg TOOLS_VERSION=$(tools_version) \
		-f deploy/local/join.results/Dockerfile ./
dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
dscan: dbuild
	docker scan --accept-license $(tag)		
#####################################################################################
clean:
	rm -rf $(dist_dir)
.PHONY:
	clean copy dbuild dpush dscan
