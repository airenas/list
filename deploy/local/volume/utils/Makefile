include Makefile.options
########################################################################################################
########################################################################################################
info: 
	@echo "==================================="
	@echo "apps dir: 	$(apps_dir)"
	@echo "apps: 		$(apps)"
	@echo "models dir: 	$(models_dir)"
	@echo "file: 		$(file)"
	@echo "URL:			$(dwn_url)"
	@echo "storage_dir:	$(storage_dir)"
	
	@echo "==================================="
########################################################################################################
$(dwn_dir):
	mkdir -p $@
$(dwn_dir)/.$(apps).done: | $(dwn_dir)
	wget --no-check-certificate -O $(dwn_dir)/$(apps).tar.gz $(dwn_url)/$(apps).tar.gz
	touch $@
$(apps_dir)/.$(apps).done: $(dwn_dir)/.$(apps).done | $(apps_dir)
	tar xvzf $(dwn_dir)/$(apps).tar.gz -C $(apps_dir)
	touch $@
$(dwn_dir)/.$(file).done: | $(dwn_dir)
	wget --no-check-certificate -O $(dwn_dir)/$(file).tar.gz $(dwn_url)/$(file).tar.gz
	touch $@	
$(models_dir)/.$(file).done: $(dwn_dir)/.$(file).done | $(models_dir)
	tar xvzf $(dwn_dir)/$(file).tar.gz -C $(models_dir)
	touch $@	
########################################################################################################
install-apps: info $(apps_dir)/.$(apps).done
########################################################################################################
install-data: info $(models_dir)/.$(file).done
########################################################################################################
$(storage_dir)/decoded/audio:
	mkdir -p $@
$(storage_dir)/decoded/trans:
	mkdir -p $@
$(storage_dir)/results:
	mkdir -p $@
$(storage_dir)/logs:
	mkdir -p $@			
init-storage: info | $(storage_dir)/decoded/audio $(storage_dir)/decoded/trans $(storage_dir)/results $(storage_dir)/logs
########################################################################################################
clean-apps:
	rm -rf $(apps_dir)
clean-data:
	rm -rf $(storage_dir)
	rm -rf $(models_dir)

PHONY: info
