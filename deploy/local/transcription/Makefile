-include ../Makefile.options
#####################################################################################
dist_dir=deploy
executable_name=cmdWorkerService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=transcription-service
src_dir=../../../src/decode
#####################################################################################
$(main_dir)/$(executable_name): $(GO_SRC_DIR)
	(cd $(main_dir) && CGO_ENABLED=0 go build -installsuffix -cgo)

build: $(main_dir)/$(executable_name)
#####################################################################################
$(dist_dir)/$(executable_name): $(main_dir)/$(executable_name)
	mkdir -p $(@D)
	cp $< $@

$(dist_dir)/cmd/Makefile: $(src_dir)/transcription/Makefile
	mkdir -p $(@D)
	cp $< $@

copy: $(dist_dir)/cmd/Makefile $(dist_dir)/$(executable_name)
########### DOCKER ##################################################################
tag=airenas/$(docker_container_prefix)$(service):$(TRANSCRIPTION_SERVICE_VERSION).$(commit_count)

dbuild: $(dist_dir)/cmd/Makefile $(dist_dir)/$(executable_name)
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
########### END #####################################################################
clean:
	rm -rf $(dist_dir)
	rm -f $(main_dir)/$(executable_name)

.PHONY:
	clean copy build dbuild dpush
