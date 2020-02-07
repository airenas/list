# make file to prepare the result from lattice
# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/decoded
RESULTS_ROOT?=/results
APPS_ROOT?=/kaldi
FINAL_PASS=2-pass

PATH:=utils:$(APPS_ROOT)/bin:$(PATH)

# init symbolic directories
steps:
	ln -s $(APPS_ROOT)/steps

utils:
	ln -s $(APPS_ROOT)/utils

.SECONDARY:
.DELETE_ON_ERROR:	

#rescore lattice
$(OUTPUT_ROOT)/trans/%/$(FINAL_PASS)/done: $(OUTPUT_ROOT)/trans/%/1-pass/decode/lat.1.gz steps utils
	mkdir -p `dirname $@`
	## run rescore script from model dir
	LM_DIR=$(LM_DIR) $(SCRIPTS_DIR)/rescore.sh $(MODELS_ROOT) $< $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode/lat.1.gz || exit 1
	touch $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/done 

$(OUTPUT_ROOT)/trans/%.selected.ctm: $(OUTPUT_ROOT)/trans/%/done
	gunzip -c $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode/lat.1.gz | lattice-1best ark:- ark:- | \
      nbest-to-ctm  ark:- - | \
      utils/int2sym.pl -f 5 $(MODELS_ROOT)/words.txt | \
	  perl -npe 's/(.*)-(S\d+)---(\S+)/\1_\3_\2/' > $@

%.synced.ctm: %.selected.ctm
	cat $^ | LC_ALL=C sort -k 1,1 -k 3,3n -k 4,4n > $@

%.txt: %.synced.ctm
	cat $^ | cut -f 5 -d " " | perl -npe 's/\n/ /'  > $@	

$(RESULTS_ROOT)/%/result.txt: $(OUTPUT_ROOT)/trans/%/$(FINAL_PASS).txt
	mkdir -p $(RESULTS_ROOT)/$*/
	cat $^ > $@	

$(RESULTS_ROOT)/%/resultFinal.txt: $(RESULTS_ROOT)/%/result.txt
ifeq ($(PUNCTUATE), true)
	$(PR_EXEC) -f $^ -u $(PR_URL) -o $@
else
	cat $^ > $@	
endif

$(RESULTS_ROOT)/%/lat.gz: $(OUTPUT_ROOT)/trans/%/$(FINAL_PASS)/done
	mkdir -p $(RESULTS_ROOT)/$*/
	cp $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode/lat.1.gz $@		

$(RESULTS_ROOT)/%/lat.txt: $(RESULTS_ROOT)/%/lat.gz
	gunzip -c $^ | $(APPS_ROOT)/bin/lattice-copy ark:- ark,t:- | \
		utils/int2sym.pl -f $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode/lat.1.gzxt: $(RESULTS_ROOT)/%/lat.gz
	gunzip -c $^ | \
	  $(APPS_ROOT)/bin/lattice-to-nbest --n=10 ark:- ark:- | \
		$(APPS_ROOT)/bin/lattice-copy ark:- ark,t:- | \
		utils/int2sym.pl -f 3 $(MODELS_ROOT)/words.txt > $@		

test: $(RESULTS_ROOT)/a/resultFinal.txt

r1: $(RESULTS_ROOT)/a/lat.gz $(RESULTS_ROOT)/a/lat.nb10.txt $(RESULTS_ROOT)/a/lat.txt $(RESULTS_ROOT)/a/resultFinal.txt