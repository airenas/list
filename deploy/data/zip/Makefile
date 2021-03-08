-include Makefile.options
MAIN_DIR=${CURDIR}/../..
OUT_DIR=$(WORK_DIR)/am/$(NAME)-$(VERSION)
OUT_NAME=$(NAME)-$(VERSION)
OUT_FILE=am-$(NAME)-$(VERSION).tar.gz
OUT_LM_FILE=lm-$(NAME)-$(VERSION).tar.gz
OUT_CONFIG_FILE=config-$(NAME)-$(VERSION).tar.gz
#################################################################################################
info:
	@echo "Input		: $(INITIAL_FILE)"
	@echo "FST		: $(FST_FILE)"
	@echo "Vocab		: $(VOCAB_FILE)"
	@echo "LM 	        : $(LM_FILE)"
	@echo "Phones		: $(PHONES_FILE)"
	@echo "Output		: $(WORK_DIR)/$(OUT_FILE)"
	@echo "                 : $(WORK_DIR)/$(OUT_CONFIG_FILE)"
	@echo "                 : $(WORK_DIR)/$(OUT_LM_FILE)"
#	@echo "Target		: $(target_files)"
#################################################################################################
$(WORK_DIR):
	mkdir -p $@
$(WORK_DIR)/initial:
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
$(WORK_DIR)/.extracted: $(INITIAL_FILE) | $(WORK_DIR)/initial
	tar xvzf $^ -C $(WORK_DIR)/initial
	rm -f $(WORK_DIR)/initial/$(INITIAL_DIR)/HCLG.fst
	touch $@
#################################################################################################
files=lmwt wip cmvn_opts final.mdl HCLG.fst phones.txt word_boundary.int words.txt \
	conf/ivector_extractor.conf conf/mfcc_hires.conf conf/vad.conf conf/online_cmvn.conf conf/decode.conf conf/splice.conf \
	ivector_extractor/splice_opts ivector_extractor/num_jobs ivector_extractor/final.ie ivector_extractor/final.dubm ivector_extractor/online_cmvn.conf \
	ivector_extractor/final.ie.id ivector_extractor/final.mat ivector_extractor/global_cmvn.stats \
	info
target_files=$(patsubst %, $(OUT_DIR)/%, $(files))

conf_files=$(OUT_NAME)/scripts/rescore.preload.sh $(OUT_NAME)/scripts/decode.preload.sh $(OUT_NAME)/scripts/decode.sh $(OUT_NAME)/scripts/rescore.sh $(OUT_NAME).yml
target_conf_files=$(patsubst %, $(WORK_DIR)/config/%, $(conf_files))

lm_files=G.fst words.txt info
target_lm_files=$(patsubst %, $(WORK_DIR)/lm/$(OUT_NAME)/%, $(lm_files))
#################################################################################################
$(FST_FILE): $(FST_FILE).gz
	pv $< | gunzip > $@
$(OUT_DIR)/HCLG.fst: $(FST_FILE) | $(OUT_DIR)
	pv $< > $@
$(OUT_DIR)/phones.txt: $(PHONES_FILE) | $(OUT_DIR)
	cp $< $@
$(OUT_DIR)/words.txt: $(VOCAB_FILE) | $(OUT_DIR)
	cp $< $@
$(OUT_DIR)/%: $(WORK_DIR)/initial/$(INITIAL_DIR)/% $(WORK_DIR)/.extracted | $(OUT_DIR)
	mkdir -p `dirname $@`
	cp $< $@
$(OUT_DIR)/info: | $(OUT_DIR)
	echo "Zip date      : $$(date -Isec)" > $@
	echo "Words         : $$(wc -l $(OUT_DIR)/words.txt | awk '{print $$1}')" >> $@
#################################################################################################
$(WORK_DIR)/lm/$(OUT_NAME)/G.fst: $(LM_FILE) | $(WORK_DIR)/lm/$(OUT_NAME)
	pv $< | gunzip > $@
$(WORK_DIR)/lm/$(OUT_NAME)/words.txt: $(VOCAB_FILE) | $(WORK_DIR)/lm/$(OUT_NAME)
	cp $< $@
$(WORK_DIR)/lm/$(OUT_NAME)/info: $(WORK_DIR)/lm/$(OUT_NAME)/words.txt
	echo "Zip date      : $$(date -Isec)" > $@
	echo "Words         : $$(wc -l $(WORK_DIR)/lm/$(OUT_NAME)/words.txt | awk '{print $$1}')" >> $@
#################################################################################################
$(WORK_DIR)/config/$(OUT_NAME)/%: $(WORK_DIR)/initial/config/words_v3/% $(WORK_DIR)/.extracted | $(WORK_DIR)/config
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
prepare: info $(WORK_DIR)/.extracted
build: info $(WORK_DIR)/.extracted $(WORK_DIR)/$(OUT_FILE)
build-config: info $(WORK_DIR)/.extracted  $(WORK_DIR)/$(OUT_CONFIG_FILE)
build-lm: info $(WORK_DIR)/$(OUT_LM_FILE)
#################################################################################################
clean:
	rm -rf $(WORK_DIR)
#################################################################################################
.PHONY: info clean
#################################################################################################
