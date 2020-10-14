include Makefile.options
########################################################################################################
install: run	
########################################################################################################
status: 
	docker-compose ps

run:
	docker-compose up -d 

clean:
	docker-compose stop && docker-compose rm -f

########################################################################################################
init-options: Makefile.options
Makefile.options: Makefile.options.template
	cp Makefile.options.template $@	
########################################################################################################
.PHONY: 