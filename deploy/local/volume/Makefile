-include ../Makefile.options

#####################################################################################
clean:
	@echo "Done"
########### DOCKER ##################################################################
tag=$(docker_prefix)volume-helper:1.1.$(commit_count)
dbuild: 
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
########### END #####################################################################
.PHONY:
	dbuild dpush
