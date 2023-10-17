#####################################################################################
dist_dir=$(CURDIR)/deploy
service=diarization-pyannote-service
version=dev
#####################################################################################
########### DOCKER ##################################################################
tag=$(service):$(version)

dbuild:
	docker buildx build -t $(tag) --build-arg BUILD_VERSION=$(version) \
		-f Dockerfile ../../

dpush: dbuild
	docker push $(tag)

dscan: dbuild
	docker scan --accept-license $(tag)	
#####################################################################################
.PHONY:
	clean copy build dbuild dpush
