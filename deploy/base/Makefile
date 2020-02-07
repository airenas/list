d_make=airenas/ubuntu.make:2.0
d_java=airenas/ubuntu.make.java8:2.0

#####################################################################################
dbuild_make: 
	docker build -t $(d_make) ubuntu.make

dpush_make: dbuild_make
	docker push $(d_make)
#####################################################################################
dbuild_java: 
	docker build -t $(d_java) ubuntu.make.java8

dpush_java: dbuild_java
	docker push $(d_java)
#####################################################################################
