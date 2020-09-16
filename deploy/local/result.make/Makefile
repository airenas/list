-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=cmdWorkerService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=result-make-service
tag=$(docker_prefix)$(service):$(TR_RESULT_SERVICE_VERSION).$(commit_count)

#####################################################################################
tools=punct.lattice lattice.to.text lattice.to.webvtt send.metric int.to.word empty.lattice \
	m_shell.sh m_start.sh m_end.sh
tools_dir=$(SRC_DIR)/tools
dist_tools_files=$(patsubst %, $(dist_dir)/%, $(tools))
#####################################################################################
script_files=Makefile restore/lat_restore.pl restore/lat_map.pl restore/lt.pl restore/LatGraph.pm \
			restore/utils_num.pl
dist_script_files=$(patsubst %, $(dist_dir)/cmd/%, $(script_files))
#####################################################################################
$(dist_dir):
	mkdir -p $@
$(dist_dir)/$(executable_name): $(GO_SRC_DIR) | $(dist_dir)
	cd $(main_dir) && $(go_build_cmd) -o $(dist_dir)/$(executable_name)
	@echo -n "$(GREEN)Build:$(C_NC) $(executable_name)\n"
$(dist_dir)/%: $(tools_dir)/% | $(dist_dir)
	cd $< && $(go_build_cmd) -o $@
	@echo -n "$(GREEN)Build:$(C_NC) $*\n"
build: $(dist_dir)/$(executable_name) $(dist_tools_files)
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
copy: $(dist_script_files)
########### DOCKER ##################################################################
dbuild: build $(dist_script_files)
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
