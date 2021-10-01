-include Makefile.options
OUT_DIR=$(WORK_DIR)/am/$(NAME)-$(VERSION)
OUT_NAME=$(NAME)-$(VERSION)
OUT_FILE=am-$(NAME)-$(VERSION).tar.gz
OUT_LM_FILE=lm-$(NAME)-$(VERSION).tar.gz
OUT_CONFIG_FILE=config-$(NAME)-$(VERSION).tar.gz
#################################################################################################
info:
	@echo "Input		: $(INITIAL_FILE)"
	@echo "Initial dir  : $(INITIAL_DIR)"
	@echo "FST			: $(FST_FILE)"
	@echo "Vocab		: $(VOCAB_FILE)"
	@echo "LM 	        : $(LM_FILE)"
	@echo "Phones		: $(PHONES_FILE)"
	@echo "Output		: $(WORK_DIR)/$(OUT_FILE)"
	@echo "             : $(WORK_DIR)/$(OUT_CONFIG_FILE)"
	@echo "             : $(WORK_DIR)/$(OUT_LM_FILE)"
#	@echo "Target		: $(target_files)"
#################################################################################################
$(WORK_DIR):
	mkdir -p $@
$(WORK_DIR)/initial:
	mkdir -p $@
$(WORK_DIR)/initial.config:
	mkdir -p $@	
$(WORK_DIR)/lm/$(OUT_NAME):
	mkdir -p $@
$(WORK_DIR)/new:
	mkdir -p $@
$(OUT_DIR):
	mkdir -p $@
$(WORK_DIR)/config:
	mkdir -p $@
########################################################################################################
$(WORK_DIR)/$(INITIAL_FILE): | $(WORK_DIR)
	wget --no-check-certificate -O $@_ $(dwn_url)/$(INITIAL_FILE)
	mv $@_ $@
	
$(WORK_DIR)/.extracted: $(WORK_DIR)/$(INITIAL_FILE) | $(WORK_DIR)/initial
	tar xvzf $^ -C $(WORK_DIR)/initial
	rm -f $(WORK_DIR)/initial/$(INITIAL_DIR)/HCLG.fst
	touch $@
########################################################################################################
$(WORK_DIR)/$(INITIAL_CONFIG_FILE): | $(WORK_DIR)
	wget --no-check-certificate -O $@_ $(dwn_url)/$(INITIAL_CONFIG_FILE)
	mv $@_ $@
	
$(WORK_DIR)/.config.extracted: $(WORK_DIR)/$(INITIAL_CONFIG_FILE) | $(WORK_DIR)/initial.config
	tar xvzf $^ -C $(WORK_DIR)/initial.config
	rm -f $(WORK_DIR)/initial/$(INITIAL_DIR)/HCLG.fst
	touch $@	
########################################################################################################
$(WORK_DIR)/$(NEW_FILE): | $(WORK_DIR)
	$(NEW_FILE_DWN_CMD)
	
$(WORK_DIR)/.new.extracted: $(WORK_DIR)/$(NEW_FILE) | $(WORK_DIR)/new
	tar xvzf $^ -C $(WORK_DIR)/new
	touch $@	
#################################################################################################
files=cmvn_opts conf/mfcc_hires.conf final.mdl frame_subsampling_factor HCLG.fst \
	ivector_extractor/online_cmvn.conf ivector_extractor/final.ie \
	ivector_extractor/final.ie.id ivector_extractor/final.dubm ivector_extractor/num_jobs \
	ivector_extractor/final.mat ivector_extractor/splice_opts ivector_extractor/global_cmvn.stats \
	lmwt phones.txt wip word_boundary.int words.txt \
	info
target_files=$(patsubst %, $(OUT_DIR)/%, $(files))

conf_files=$(OUT_NAME)/scripts/rescore.preload.sh $(OUT_NAME)/scripts/decode.preload.sh $(OUT_NAME)/scripts/decode.sh $(OUT_NAME)/scripts/rescore.sh $(OUT_NAME).yml
target_conf_files=$(patsubst %, $(WORK_DIR)/config/%, $(conf_files))

lm_files=G.fst words.txt info
target_lm_files=$(patsubst %, $(WORK_DIR)/lm/$(OUT_NAME)/%, $(lm_files))
#################################################################################################
$(OUT_DIR)/%: $(WORK_DIR)/new/$(NEW_DIR)/% $(WORK_DIR)/.new.extracted | $(OUT_DIR)
	mkdir -p `dirname $@`
	cp $< $@	
$(OUT_DIR)/%: $(WORK_DIR)/initial/$(INITIAL_DIR)/% $(WORK_DIR)/.extracted | $(OUT_DIR)
	mkdir -p `dirname $@`
	cp $< $@
$(OUT_DIR)/info: | $(OUT_DIR)
	echo "Build date    : $$(date -Isec)" > $@
	echo "Words         : $$(wc -l $(OUT_DIR)/words.txt | awk '{print $$1}')" >> $@
	echo "From base     : $(INITIAL_FILE)" >> $@
	echo "From          : $(NEW_FILE)" >> $@
#################################################################################################
$(WORK_DIR)/lm/$(OUT_NAME)/G.fst: $(WORK_DIR)/new/$(NEW_DIR)/G.fst | $(WORK_DIR)/lm/$(OUT_NAME)
	cp $< $@
$(WORK_DIR)/lm/$(OUT_NAME)/words.txt: $(WORK_DIR)/new/$(NEW_DIR)/words.txt | $(WORK_DIR)/lm/$(OUT_NAME)
	cp $< $@
$(WORK_DIR)/lm/$(OUT_NAME)/info: $(WORK_DIR)/lm/$(OUT_NAME)/words.txt
	echo "Build date    : $$(date -Isec)" > $@
	echo "Words         : $$(wc -l $(WORK_DIR)/lm/$(OUT_NAME)/words.txt | awk '{print $$1}')" >> $@
	echo "From          : $(NEW_FILE)" >> $@
#################################################################################################
$(WORK_DIR)/config/$(OUT_NAME)/%: $(WORK_DIR)/initial.config/config/$(INITIAL)/% $(WORK_DIR)/.extracted | $(WORK_DIR)/config
	mkdir -p `dirname $@`
	cp $< $@
$(WORK_DIR)/config/$(OUT_NAME).yml: | $(WORK_DIR)/config
	echo 'name: $(OUT_NAME)' > $@
	echo 'description: Žodžių atpažintuvas - $(OUT_NAME).' >> $@
	echo "date_created: $$(date -Idate)" >> $@
	echo 'settings:' >> $@
	echo '  models_root: /models/am/$(OUT_NAME)' >> $@
	echo '  scripts_dir: /models/config/$(OUT_NAME)/scripts/' >> $@
	echo '  lm_dir: /models/lm/$(OUT_NAME)/' >> $@
	echo '  rnnlm_dir: /models/rnnlm/2.0/' >> $@
	echo '  punctuate: true' >> $@
	echo '  transcription-preload_key: $(OUT_NAME)' >> $@
	echo '  transcription-preload_cmd: /models/config/$(OUT_NAME)/scripts/decode.preload.sh' >> $@
	echo '  rescore-preload_key: $(OUT_NAME)' >> $@
	echo '  rescore-preload_cmd: /models/config/$(OUT_NAME)/scripts/rescore.preload.sh' >> $@
#################################################################################################
$(WORK_DIR)/$(OUT_FILE): $(target_files)
	cd $(WORK_DIR) && tar -czvf $(OUT_FILE) am
#################################################################################################
$(WORK_DIR)/$(OUT_CONFIG_FILE): $(target_conf_files)
	cd $(WORK_DIR) && tar -czvf $(OUT_CONFIG_FILE) config
#################################################################################################
$(WORK_DIR)/$(OUT_LM_FILE): $(target_lm_files)
	cd $(WORK_DIR) && tar -czvf $(OUT_LM_FILE) lm
#################################################################################################
prepare: info $(WORK_DIR)/.extracted $(WORK_DIR)/.new.extracted
build: info $(WORK_DIR)/.extracted $(WORK_DIR)/.new.extracted $(WORK_DIR)/$(OUT_FILE)
build-config: info $(WORK_DIR)/.config.extracted $(WORK_DIR)/$(OUT_CONFIG_FILE)
build-lm: info $(WORK_DIR)/$(OUT_LM_FILE)
test: $(WORK_DIR)/.extracted $(WORK_DIR)/.new.extracted
#################################################################################################
SSH_CMD=ssh -t -p $(SSH_PORT) $(USER)@$(SERVER_URL)
RSYNC_CMD=rsync -r -e 'ssh -p $(SSH_PORT)' -P
fake-%.upload: $(WORK_DIR)/%.tar.gz
	$(RSYNC_CMD) $< $(USER)@$(SERVER_URL):$(remote_dwn_dir)/
install-%.data: fake-%.upload
	$(SSH_CMD) "tar xvzf $(remote_dwn_dir)/$*.tar.gz -C $(remote_models_dir)/"

install-am: install-am-$(OUT_NAME).data	
install-config: install-config-$(OUT_NAME).data	
install-lm: install-lm-$(OUT_NAME).data	

backup-%.upload: $(WORK_DIR)/%.tar.gz
	$(RSYNC_CMD) $< $(BACKUP_URL)/
backup-am: backup-am-$(OUT_NAME).upload	
backup-config: backup-config-$(OUT_NAME).upload	
backup-lm: backup-lm-$(OUT_NAME).upload	
#################################################################################################
clean:
	rm -rf $(WORK_DIR)
#################################################################################################
.PHONY: info clean
#################################################################################################
