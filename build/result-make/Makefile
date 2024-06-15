#####################################################################################
dist_dir=$(CURDIR)/deploy
service=result-make-service
version=dev
#####################################################################################
########### DOCKER ##################################################################
tag=airenas/list-$(service):$(version)

dbuild:
	docker buildx build -t $(tag) --build-arg BUILD_VERSION=$(version) --build-arg TOOLS_VERSION=$(version) \
		-f Dockerfile ../../

dpush: dbuild
	docker push $(tag)

dscan: dbuild
	docker scan --accept-license $(tag)	
#####################################################################################
.PHONY:
	clean copy build dbuild dpush
