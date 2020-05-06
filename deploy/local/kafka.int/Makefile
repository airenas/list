-include ../Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/deploy
executable_name=kafkaIntegrationService
main_dir=$(GO_SRC_DIR)/cmd/$(executable_name)
service=kafka-integration-service
version=$(KAFKA_INTEGRATION_SERVICE_VERSION)
#####################################################################################
$(dist_dir)/$(executable_name): $(GO_SRC_DIR)
	cd $(main_dir) && go build -tags static -o $@

build: $(dist_dir)/$(executable_name)
#####################################################################################
tag=airenas/$(docker_container_prefix)$(service):$(version).$(commit_count)

dbuild: $(dist_dir)/$(executable_name)
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
########### END #####################################################################
clean:
	rm -rf $(dist_dir)
	rm -f $(main_dir)/$(executable_name)

.PHONY:
	clean build dbuild dpush

