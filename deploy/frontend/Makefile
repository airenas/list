include Makefile.options
########################################################################################################
install: run	
########################################################################################################
status: 
	docker-compose ps

run: traefik.toml rules.toml
	docker-compose up -d 

clean:
	docker-compose stop && docker-compose rm -f

restart: traefik.toml rules.toml
	docker-compose stop $(srv) && docker-compose rm -f $(srv)
	docker-compose up -d $(srv)
 
########################################################################################################
init-options: Makefile.options
Makefile.options: Makefile.options.template
	cp Makefile.options.template $@	
traefik.toml: traefik.toml.in Makefile.options
	cat $< | envsubst > $@
rules.toml: rules.toml.in Makefile.options
	cat $< | envsubst > $@	

########################################################################################################
.PHONY: 
.EXPORT_ALL_VARIABLES:
