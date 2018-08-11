# make file to prepare the result from lattice
# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/decoded
APPS_ROOT?=/kaldi
LM_SCALE?=17
FINAL_PASS?=1-pass

PATH:=utils:$(APPS_ROOT)/bin:$(PATH)

# init symbolic directories
steps:
	ln -s $(APPS_ROOT)/steps

utils:
	ln -s $(APPS_ROOT)/utils

.SECONDARY:
.DELETE_ON_ERROR:	

#extract data from lattice
%/decode/.ctm: %/decode/log steps utils
	frame_shift_opt=""; \
	if [ -f $*/frame_subsampling_factor ]; then \
	  factor=`cat $*/frame_subsampling_factor`; \
	  frame_shift_opt="--frame-shift 0.0$$factor"; \
	fi; \
	steps/get_ctm.sh $$frame_shift_opt `dirname $*` $* $*/decode
	touch -m $@

$(OUTPUT_ROOT)/trans/%.selected.ctm: $(OUTPUT_ROOT)/trans/%/decode/.ctm
	cat $(OUTPUT_ROOT)/trans/$*/decode/score_$(LM_SCALE)/`dirname $*`.ctm  | perl -npe 's/(.*)-(S\d+)---(\S+)/\1_\3_\2/' > $@

%.txt: %.selected.ctm
	cat $^ | cut -f 5 -d " " | perl -npe 's/\n/ /'  > $@

$(OUTPUT_ROOT)/trans/%/result.txt: $(OUTPUT_ROOT)/trans/%/$(FINAL_PASS).txt
	cat $^ > $@	

