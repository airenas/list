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

###############################################################################
SEND_METRIC_APP?=${APP_DIR}/send.metric
###############################################################################

# init symbolic directories
steps:
	ln -s $(APPS_ROOT)/steps

utils:
	ln -s $(APPS_ROOT)/utils

.SECONDARY:
.DELETE_ON_ERROR:

###############################################################################
## for metrics
SHELL=$(APP_DIR)/m_shell.sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
worker=transcription
###############################################################################

info:
	echo "=====================INFO=========================="
	echo "OUTPUT_ROOT=$(OUTPUT_ROOT)"
	echo "INPUT_ROOT= $(INPUT_ROOT)"
	echo "SCRIPTS_DIR=$(SCRIPTS_DIR)"
	echo "MODELS_ROOT=$(MODELS_ROOT)"
	echo "APPS_ROOT=  $(APPS_ROOT)"
	echo "==================================================="

# MFCC calculation
$(OUTPUT_ROOT)/trans/%/mfcc/.done: info $(INPUT_ROOT)/trans/%/spk2utt steps utils
	tsk=mfcc;id=$*;. $(APP_DIR)/m_start.sh $@

	rm -rf $(OUTPUT_ROOT)/trans/$*/mfcc
	rm -f $(OUTPUT_ROOT)/trans/$*/cmvn.scp
	steps/make_mfcc.sh --mfcc-config $(MODELS_ROOT)/conf/mfcc_hires.conf --cmd "$$train_cmd" --nj $(njobs) \
		$(INPUT_ROOT)/trans/$* $(OUTPUT_ROOT)/trans/$*/exp/make_mfcc $(OUTPUT_ROOT)/trans/$*/mfcc || exit 1
	steps/compute_cmvn_stats.sh $(OUTPUT_ROOT)/trans/$* $(OUTPUT_ROOT)/trans/$*/exp/make_mfcc $(OUTPUT_ROOT)/trans/$*/mfcc || exit 1
	touch $@

$(OUTPUT_ROOT)/trans/%/1-pass/.ivectors.done: $(OUTPUT_ROOT)/trans/%/mfcc/.done
	tsk=extractIVectors;id=$*;. $(APP_DIR)/m_start.sh $@

	rm -rf $(OUTPUT_ROOT)/trans/$*/1-pass
	mkdir -p $(OUTPUT_ROOT)/trans/$*/1-pass/decode
	steps/online/nnet2/extract_ivectors_online.sh --cmd "$$decode_cmd" --nj $(njobs) \
        $(OUTPUT_ROOT)/trans/$* $(MODELS_ROOT)/ivector_extractor $(OUTPUT_ROOT)/trans/$*/nnet2_online/ivectors
	(cd $(OUTPUT_ROOT)/trans/$*/1-pass; for f in $(MODELS_ROOT)/*; do ln -s $$f; done)
	touch $@	

### Do 1-pass decoding using chain online models
$(OUTPUT_ROOT)/trans/%/1-pass/.done: info $(OUTPUT_ROOT)/trans/%/1-pass/.ivectors.done
	tsk=transcribe;id=$*;. $(APP_DIR)/m_start.sh $@
	## run decode script from model dir
	echo Scripts dir=$(SCRIPTS_DIR)
	$(SCRIPTS_DIR)/decode.sh $(MODELS_ROOT) \
		$(OUTPUT_ROOT)/trans/$*/nnet2_online/ivectors \
		$(OUTPUT_ROOT)/trans/$* $(OUTPUT_ROOT)/trans/$*/1-pass/decode || exit 1
	touch $@	

.PHONY:
	info