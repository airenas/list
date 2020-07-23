-include Makefile.options
MAIN_DIR=${CURDIR}/../..
OUT_DIR=$(WORK_DIR)/am/$(NAME)-$(VERSION)
OUT_FILE=am-$(NAME)$-(VERSION).tar.gz
#################################################################################################
info:
	@echo "Input		: $(INITIAL_FILE)"
	@echo "FST		: $(FST_FILE)"
	@echo "Vocab		: $(VOCAB_FILE)"
	@echo "Phones		: $(PHONES_FILE)"
	@echo "Output		: $(WORK_DIR)/$(OUT_FILE)"
	@echo "TrainStep 	: $(STEP)"
#################################################################################################
$(WORK_DIR):
	mkdir -p $@
$(WORK_DIR)/initial:
	mkdir -p $@
$(WORK_DIR)/new:
	mkdir -p $@
$(OUT_DIR):
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
#################################################################################################
$(OUT_DIR)/HCLG.fst: $(FST_FILE) | $(OUT_DIR)
	cat $< | gunzip > $@
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
$(WORK_DIR)/$(OUT_FILE): $(target_files)
	cd $(WORK_DIR) && tar -czvf $(OUT_FILE) rnnlm
#################################################################################################
build: info $(WORK_DIR)/.extracted $(target_files)
#################################################################################################
clean:
	rm -rf $(WORK_DIR)
#################################################################################################
.PHONY: info clean
#################################################################################################
