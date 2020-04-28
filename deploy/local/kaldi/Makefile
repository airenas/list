-include Makefile.options
#####################################################################################
RSYNC_CMD=rsync -r -P 
#####################################################################################
dist_dir=apps
kaldi_bin=$(dist_dir)/kaldi/bin
kaldi_steps=$(dist_dir)/kaldi/steps
kaldi_utils=$(dist_dir)/kaldi/utils
#####################################################################################
bin_files=apply-cmvn arpa-to-const-arpa compute-cmvn-stats compute-mfcc-feats compute-wer \
	copy-feats ivector-extract-online2 lattice-1best lattice-add-penalty lattice-align-words-lexicon \
	lattice-best-path lattice-copy lattice-lmrescore lattice-lmrescore-const-arpa lattice-minimize \
	lattice-prune lattice-scale lattice-to-nbest nbest-to-ctm nnet3-latgen-faster nnet3-latgen-faster-parallel \
	extract-segments lattice-push \
	nnet3-latgen-faster-parallel-pipe \
	lattice-lmrescore-kaldi-rnnlm-pruned rnnlm-get-word-embedding \
	lattice-lmrescore-kaldi-rnnlm-pruned-pipe
#####################################################################################
d_bin_files=$(patsubst %, $(dist_dir)/kaldi/bin/%, $(bin_files))
#####################################################################################
$(dist_dir)/kaldi: 
	mkdir -p $@
$(dist_dir)/kaldi/bin: 
	mkdir -p $@
########### BIN ###################################################################
$(dist_dir)/kaldi/bin/%: $(KALDI_ROOT)/src/bin/% | $(dist_dir)/kaldi/bin 
	cp $< $@
$(dist_dir)/kaldi/bin/%: $(KALDI_ROOT)/src/featbin/%  | $(dist_dir)/kaldi/bin
	cp $< $@	
$(dist_dir)/kaldi/bin/%: $(KALDI_ROOT)/src/lmbin/%  | $(dist_dir)/kaldi/bin
	cp $< $@
$(dist_dir)/kaldi/bin/%: $(KALDI_ROOT)/src/online2bin/%  | $(dist_dir)/kaldi/bin
	cp $< $@
$(dist_dir)/kaldi/bin/%: $(KALDI_ROOT)/src/latbin/%  | $(dist_dir)/kaldi/bin
	cp $< $@			
$(dist_dir)/kaldi/bin/%: $(KALDI_ROOT)/src/nnet3bin/%  | $(dist_dir)/kaldi/bin
	cp $< $@
$(dist_dir)/kaldi/bin/%: $(KALDI_ROOT)/src/rnnlmbin/%  | $(dist_dir)/kaldi/bin
	cp $< $@				
build:
	cd $(KALDI_ROOT)/$(kd) && $(MAKE) $(tool)
########### STEPS ###################################################################
$(dist_dir)/kaldi/steps: $(dist_dir)/kaldi
	$(RSYNC_CMD) $(KALDI_ROOT)/egs/wsj/s5/steps $(dist_dir)/kaldi/

$(dist_dir)/kaldi/utils: $(dist_dir)/kaldi
	$(RSYNC_CMD) $(KALDI_ROOT)/egs/wsj/s5/utils $(dist_dir)/kaldi/	

$(dist_dir)/kaldi/utils/run.pl: $(dist_dir)/kaldi/utils/parallel/run.pl
	cp $< $@	

copy: $(dist_dir)/kaldi/utils $(dist_dir)/kaldi/steps $(d_bin_files) $(dist_dir)/kaldi/utils/run.pl
########### END #####################################################################
clean:
	rm -rf $(dist_dir)

.PHONY:
	copy
