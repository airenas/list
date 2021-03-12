-include Makefile.options

OUTPUT_ROOT?=/decoded
APPS_ROOT?=/kaldi
njobs=1
nthreads=4
cdir=$(shell pwd)

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

###############################################################################
.ONESHELL:
.EXPORT_ALL_VARIABLES:
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
$(TRANS_DIR)/%/mfcc/.done: info $(TRANS_DIR)/%/spk2utt steps utils

	rm -rf $(TRANS_DIR)/$*/mfcc
	rm -f $(TRANS_DIR)/$*/cmvn.scp
	steps/make_mfcc.sh --mfcc-config $(MODELS_ROOT)/conf/mfcc_hires.conf --cmd "$$train_cmd" --nj $(njobs) \
		$(TRANS_DIR)/$* $(TRANS_DIR)/$*/exp/make_mfcc $(TRANS_DIR)/$*/mfcc || exit 1
	steps/compute_cmvn_stats.sh $(TRANS_DIR)/$* $(TRANS_DIR)/$*/exp/make_mfcc $(TRANS_DIR)/$*/mfcc || exit 1
	touch $@

$(TRANS_DIR)/%/1-pass/.ivectors.done: $(TRANS_DIR)/%/mfcc/.done

	rm -rf $(OUTPUT_ROOT)/trans/$*/1-pass
	mkdir -p $(OUTPUT_ROOT)/trans/$*/1-pass/decode
	steps/online/nnet2/extract_ivectors_online.sh --cmd "$$decode_cmd" --nj $(njobs) \
        $(TRANS_DIR)/$* $(MODELS_ROOT)/ivector_extractor $(TRANS_DIR)/$*/nnet2_online/ivectors
	(cd $(TRANS_DIR)/$*/1-pass; for f in $(MODELS_ROOT)/*; do ln -s $$f; done)
	touch $@

### Do 1-pass decoding using chain online models
$(OUTPUT_ROOT)/trans/%/1-pass/.done: info $(OUTPUT_ROOT)/trans/%/1-pass/.ivectors.done
	nnet3-latgen-faster \
		--online-ivectors=scp:$(TRANS_DIR)/$*/nnet2_online/ivectors/ivector_online.scp \
		--online-ivector-period=10 \
		--frame-subsampling-factor=3 \
		--frames-per-chunk=90 \
		--extra-left-context=0 \
		--extra-right-context=0 \
		--extra-left-context-initial=-1 \
		--extra-right-context-final=-1 \
		--minimize=false --max-active=7000 --min-active=200 --beam=15 \
		--lattice-beam=8 --acoustic-scale=1.0 --allow-partial=true \
		$(MODELS_ROOT)/final.mdl $(MODELS_ROOT)/HCLG.fst \
		"ark,s,cs:apply-cmvn --norm-means=false --norm-vars=false --utt2spk=ark:$(TRANS_DIR)/$*/utt2spk scp:$(TRANS_DIR)/$*/cmvn.scp scp:$(TRANS_DIR)/$*/feats.scp ark:- |" \
		"ark:|lattice-scale --acoustic-scale=10.0 ark:- ark:- | gzip -c >$(TRANS_DIR)/$*/1-pass/decode/lat.1.gz" || exit 1;
	touch $@

$(OUTPUT_ROOT)/trans/%/2-pass/.done: $(OUTPUT_ROOT)/trans/%/1-pass/.done
	lattice-scale --lm-scale=12.0 "ark:gunzip -c $(TRANS_DIR)/$*/1-pass/decode/lat.1.gz|" "ark:| gzip -c > $(TRANS_DIR)/$*/2-pass/decode/lat.1.gz" || exit 1
	touch $@


build: $(TRANS_DIR)/$(ID)/2-pass/.done

clean:
	rm -rf $(TRANS_DIR)/$(ID)/mfcc
	rm -f $(TRANS_DIR)/$(ID)/cmvn.scp
	rm -f steps
	rm -f utils
	rm -rf $(TRANS_DIR)/$(ID)/2-pass
	rm -rf $(TRANS_DIR)/$(ID)/1-pass

start-docker:
	docker run -it --rm --entrypoint '/bin/bash' -v $(cdir):/app \
		-v $(LOCAL_VOLUME_ROOT):$(VOLUME_ROOT) -v $(LOCAL_VOLUME_ROOT)/fileStorage/decoded/audio:/data/decoded/audio \
		airenas/list-transcription-service:3.2.633 -c '/bin/bash'
.PHONY:
	info
