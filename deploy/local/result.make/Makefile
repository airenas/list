-include ../Makefile.options
#####################################################################################
dist_dir=deploy
executable_name=cmdWorkerService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=result-make-service
tag=airenas/$(docker_container_prefix)$(service):$(TR_RESULT_SERVICE_VERSION).$(commit_count)

punct_wrapper_name=punct.file
punct_wrapper_dir=$(SRC_DIR)/tools/$(punct_wrapper_name)
#####################################################################################
$(main_dir)/$(executable_name): $(GO_SRC_DIR)
	(cd $(main_dir) && $(go_build_cmd))
$(punct_wrapper_dir)/$(punct_wrapper_name): 
	(cd $(punct_wrapper_dir) && $(go_build_cmd))

build: $(main_dir)/$(executable_name) $(punct_wrapper_dir)/$(punct_wrapper_name)
#####################################################################################
$(dist_dir)/$(executable_name): $(main_dir)/$(executable_name)
	mkdir -p $(@D)
	cp $< $@

$(dist_dir)/$(punct_wrapper_name): $(punct_wrapper_dir)/$(punct_wrapper_name)
	mkdir -p $(@D)
	cp $< $@	

$(dist_dir)/cmd/Makefile: $(SRC_DIR)/decode/result.make/Makefile
	mkdir -p $(@D)
	cp $< $@

copy: $(dist_dir)/cmd/Makefile $(dist_dir)/$(executable_name) $(dist_dir)/$(punct_wrapper_name)
########### DOCKER ##################################################################
dbuild: copy
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
	@echo -n "\n$(GREEN)CONTAINER:$(C_NC) $(tag)\n\n"
########### END #####################################################################
clean:
	rm -rf $(dist_dir)
	rm -f $(main_dir)/$(executable_name) $(punct_wrapper_dir)/$(punct_wrapper_name)

.PHONY:
	clean copy build dbuild dpush
