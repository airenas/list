# make file to transcribe audio with kaldi
# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/decoded
APPS_ROOT?=/kaldi
njobs=1
nthreads=4

PATH:=utils:$(APPS_ROOT)/bin:$(PATH)
export train_cmd=run.pl
export decode_cmd=run.pl

# init symbolic directories
steps:
	ln -s $(APPS_ROOT)/steps

utils:
	ln -s $(APPS_ROOT)/utils

.SECONDARY:
.DELETE_ON_ERROR:

log=@echo $$(date +'%Y-%m-%d %T')
logb=$(log): START
loge=$(log): END

info:
	@echo "=====================INFO=========================="
	@echo "OUTPUT_ROOT=$(OUTPUT_ROOT)"
	@echo "INPUT_ROOT= $(INPUT_ROOT)"
	@echo "SCRIPTS_DIR=$(SCRIPTS_DIR)"
	@echo "MODELS_ROOT=$(MODELS_ROOT)"
	@echo "APPS_ROOT=  $(APPS_ROOT)"
	@echo "==================================================="

# MFCC calculation
$(OUTPUT_ROOT)/trans/%/mfcc/.done: info $(INPUT_ROOT)/trans/%/spk2utt steps utils
	$(logb) $@
	rm -rf $(OUTPUT_ROOT)/trans/$*/mfcc
	rm -f $(OUTPUT_ROOT)/trans/$*/cmvn.scp
	steps/make_mfcc.sh --mfcc-config $(MODELS_ROOT)/conf/mfcc_hires.conf --cmd "$$train_cmd" --nj $(njobs) \
		$(INPUT_ROOT)/trans/$* $(OUTPUT_ROOT)/trans/$*/exp/make_mfcc $(OUTPUT_ROOT)/trans/$*/mfcc || exit 1
	steps/compute_cmvn_stats.sh $(OUTPUT_ROOT)/trans/$* $(OUTPUT_ROOT)/trans/$*/exp/make_mfcc $(OUTPUT_ROOT)/trans/$*/mfcc || exit 1
	touch $@
	$(loge) $@

### Do 1-pass decoding using chain online models
$(OUTPUT_ROOT)/trans/%/1-pass/.done: info $(OUTPUT_ROOT)/trans/%/mfcc/.done
	$(logb) $@
	rm -rf $(OUTPUT_ROOT)/trans/$*/1-pass
	mkdir -p $(OUTPUT_ROOT)/trans/$*/1-pass/decode
	steps/online/nnet2/extract_ivectors_online.sh --cmd "$$decode_cmd" --nj $(njobs) \
        $(OUTPUT_ROOT)/trans/$* $(MODELS_ROOT)/ivector_extractor $(OUTPUT_ROOT)/trans/$*/nnet2_online/ivectors
	(cd $(OUTPUT_ROOT)/trans/$*/1-pass; for f in $(MODELS_ROOT)/*; do ln -s $$f; done)
	$(logb) $@ decode
	## run decode script from model dir
	echo Scripts dir=$(SCRIPTS_DIR)
	$(SCRIPTS_DIR)/decode.sh $(MODELS_ROOT) \
		$(OUTPUT_ROOT)/trans/$*/nnet2_online/ivectors \
		$(OUTPUT_ROOT)/trans/$* $(OUTPUT_ROOT)/trans/$*/1-pass/decode || exit 1

	(cd $(OUTPUT_ROOT)/trans/$*/1-pass; ln -s ../../../fst/model)
	touch $@	
	$(loge) $@

.PHONY:
	info