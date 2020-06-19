-include ../Makefile.options
########### DOCKER ##################################################################
tag=airenas/$(docker_container_prefix)volume-helper:1.1.$(commit_count)
dbuild: 
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
########### END #####################################################################
.PHONY:
	dbuild dpush
