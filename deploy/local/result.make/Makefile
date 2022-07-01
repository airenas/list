-include ../Makefile.options
-include ../../../version
#####################################################################################
dist_dir=$(CURDIR)/deploy
service=result-make-service
tag=$(docker_prefix)$(service):$(result_make_version)
#####################################################################################
tools=m_shell.sh m_start.sh m_end.sh
tools_dir=$(SRC_DIR)/tools
dist_tools_files=$(patsubst %, $(dist_dir)/%, $(tools))
#####################################################################################
#####################################################################################
script_files=Makefile restore/lat_restore.pl restore/lat_map.pl restore/lt.pl restore/LatGraph.pm \
			restore/utils_num.pl
dist_script_files=$(patsubst %, $(dist_dir)/cmd/%, $(script_files))
#####################################################################################
$(dist_dir):
	mkdir -p $@
#####################################################################################
$(dist_dir)/cmd: 
	mkdir -p $@
$(dist_dir)/cmd/restore: 
	mkdir -p $@
$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/result.make/Makefile | $(dist_dir)/cmd
	cp $< $@
$(dist_dir)/%.sh: $(SRC_DIR)/decode/scripts/%.sh | $(dist_dir)
	cp $< $@
$(dist_dir)/cmd/restore/%: $(SRC_DIR)/decode/result.make/restore/% | $(dist_dir)/cmd/restore
	cp $< $@	
########### DOCKER ##################################################################
dbuild: $(dist_script_files) $(dist_tools_files)
	cd ../../../ && docker build -t $(tag) --build-arg BUILD_VERSION=$(version) --build-arg TOOLS_VERSION=$(tools_version) \
		-f deploy/local/result.make/Dockerfile ./
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
#####################################################################################
