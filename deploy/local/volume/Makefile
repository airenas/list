-include ../Makefile.options
########### DOCKER ##################################################################
tag=airenas/volume-helper:1.0
dbuild: 
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
########### END #####################################################################
.PHONY:
	dbuild dpush
