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

conf:
	ln -s $(MODELS_ROOT)/conf

# MFCC calculation
$(OUTPUT_ROOT)/trans/%/mfcc: $(INPUT_ROOT)/trans/%/spk2utt steps utils conf
	rm -rf $@
	rm -f $(OUTPUT_ROOT)/trans/$*/cmvn.scp
	steps/make_mfcc.sh --mfcc-config conf/mfcc_hires.conf --cmd "$$train_cmd" --nj $(njobs) \
		$(INPUT_ROOT)/trans/$* $(OUTPUT_ROOT)/trans/$*/exp/make_mfcc $@ || exit 1
	steps/compute_cmvn_stats.sh $(OUTPUT_ROOT)/trans/$* $(OUTPUT_ROOT)/trans/$*/exp/make_mfcc $@ || exit 1


### Do 1-pass decoding using chain online models
$(OUTPUT_ROOT)/trans/%/1-pass/decode/log: $(OUTPUT_ROOT)/trans/%/mfcc
	rm -rf $(OUTPUT_ROOT)/trans/$*/1-pass
	mkdir -p $(OUTPUT_ROOT)/trans/$*/1-pass
	steps/online/nnet2/extract_ivectors_online.sh --cmd "$$decode_cmd" --nj $(njobs) \
        $(OUTPUT_ROOT)/trans/$* $(MODELS_ROOT)/ivector_extractor $(OUTPUT_ROOT)/trans/$*/nnet2_online/ivectors
	(cd $(OUTPUT_ROOT)/trans/$*/1-pass; for f in $(MODELS_ROOT)/*; do ln -s $$f; done)
	steps/nnet3/decode.sh --num-threads $(nthreads) --acwt 1.0  --post-decode-acwt 10.0 \
	    --config conf/decode.conf --skip-scoring true --cmd "$$decode_cmd" --nj $(njobs) \
	    --online-ivector-dir $(OUTPUT_ROOT)/trans/$*/nnet2_online/ivectors \
	    --skip-diagnostics true \
      $(MODELS_ROOT) $(OUTPUT_ROOT)/trans/$* `dirname $@` || exit 1;

	(cd $(OUTPUT_ROOT)/trans/$*/1-pass; ln -s ../../../fst/model)