-include Makefile.options
MAIN_DIR=${CURDIR}/../..
OUT_DIR=$(WORK_DIR)/am/$(NAME)-$(VERSION)
OUT_FILE=am-$(NAME)$-(VERSION).tar.gz
#################################################################################################
info:
	@echo "Input		: $(INITIAL_FILE)"
	@echo "FST			: $(FST_FILE)"
	@echo "Vocab		: $(VOCAB_FILE)"
	@echo "Phones		: $(PHONES_FILE)"
	@echo "Output		: $(WORK_DIR)/$(OUT_FILE)"
	@echo "TrainStep 	: $(STEP)"
#################################################################################################
$(WORK_DIR):
	mkdir -p $@
$(WORK_DIR)/initial:
	mkdir -p $@
$(OUT_DIR):
	mkdir -p $@
########################################################################################################
$(WORK_DIR)/.extracted: $(INITIAL_FILE) | $(WORK_DIR)/initial
	tar xvzf $^ -C $(WORK_DIR)/initial
	touch $@	
#################################################################################################
files=feat_embedding.final.mat info.txt word_feats.txt final.raw special_symbol_opts.txt info
target_files=$(patsubst %, $(OUT_DIR)/%, $(files))
#################################################################################################
$(OUT_DIR)/feat_embedding.final.mat: | $(OUT_DIR)
	cp $(INPUT_DIR)/lm/feat_embedding.$(STEP).mat $@
$(OUT_DIR)/final.raw: | $(OUT_DIR)
	cp $(INPUT_DIR)/lm/$(STEP).raw $@
$(OUT_DIR)/%.txt: $(INPUT_DIR)/lm/%.txt | $(OUT_DIR)
	cp $^ $@
$(OUT_DIR)/info: | $(OUT_DIR)
	echo "Training step : $(STEP)" > $@
	echo "Zip date      : $$(date -Isec)" >> $@
	echo "Words         : $$(wc -l $(INPUT_DIR)/lm/config/words.txt | awk '{print $$1}')" >> $@
	echo "Features      : $$(wc -l $(INPUT_DIR)/lm/config/features.txt | awk '{print $$1}')" >> $@
	echo "Train on words: $$(cat $(INPUT_DIR)/text/lt.counts | awk '{s+=$$2}END{print s}')" >> $@
#################################################################################################
$(WORK_DIR)/$(OUT_FILE): $(target_files)
	cd $(WORK_DIR) && tar -czvf $(OUT_FILE) rnnlm
#################################################################################################
build: info $(WORK_DIR)/.extracted	
#################################################################################################
clean:
	rm -rf $(WORK_DIR)
#################################################################################################
.PHONY: info clean
#################################################################################################
