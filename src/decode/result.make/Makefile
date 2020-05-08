# make file to prepare the result from lattice
# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/decoded
RESULTS_ROOT?=/results
APPS_ROOT?=/kaldi
FINAL_PASS=2-pass
APP_DIR?=/app

PATH:=utils:$(APPS_ROOT)/bin:$(PATH)

log=echo $$(date +'%Y-%m-%d %T')
logb=$(log): START
loge=$(log): END

###############################################################################
## for metrics
SHELL=$(APP_DIR)/m_shell.sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
worker=result.make
###############################################################################
# init symbolic directories
steps:
	ln -s $(APPS_ROOT)/steps

utils:
	ln -s $(APPS_ROOT)/utils

.SECONDARY:
.DELETE_ON_ERROR:

#scale
$(OUTPUT_ROOT)/trans/%/restore/lat.scaled.gz: $(OUTPUT_ROOT)/trans/%/$(FINAL_PASS)/.done
	tsk=scale;id=$*;. $(APP_DIR)/m_start.sh $@

	mkdir -p $(OUTPUT_ROOT)/trans/$*/restore
	opt_lmwt=$$(cat "$(MODELS_ROOT)/lmwt"); \
	opt_wip=$$(cat "$(MODELS_ROOT)/wip"); \
	lattice-scale --inv-acoustic-scale=$${opt_lmwt} \
		"ark:gunzip -c $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode/lat.1.gz|"  ark:- | \
		lattice-add-penalty --word-ins-penalty=$${opt_wip} ark:- "ark:| gzip -c > $@"

$(RESULTS_ROOT)/%/result.txt: $(OUTPUT_ROOT)/trans/%/restore/.done
	tsk=lattice2text;id=$*;. $(APP_DIR)/m_start.sh $@
	
	mkdir -p $(RESULTS_ROOT)/$*/
	$(APP_DIR)/lattice.to.text $(OUTPUT_ROOT)/trans/$*/restore/lat.restored.txt $@

$(RESULTS_ROOT)/%/resultFinal.txt: $(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt
	tsk=lattice2textFinal;id=$*;. $(APP_DIR)/m_start.sh $@

	mkdir -p $(RESULTS_ROOT)/$*/
	$(APP_DIR)/lattice.to.text $^ $@

$(RESULTS_ROOT)/%/webvtt.txt: $(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt
	tsk=lattice2webvtt;id=$*;. $(APP_DIR)/m_start.sh $@
	
	mkdir -p $(RESULTS_ROOT)/$*/
	$(APP_DIR)/lattice.to.webvtt $^ $@

$(RESULTS_ROOT)/%/lat.restored.txt: $(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt
	mkdir -p $(RESULTS_ROOT)/$*/
	cp $^ $@

$(RESULTS_ROOT)/%/lat.restored.gz: $(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt
	mkdir -p $(RESULTS_ROOT)/$*/
	gzip -c $^ > $@

$(RESULTS_ROOT)/%/lat.gz: $(OUTPUT_ROOT)/trans/%/restore/lat.scaled.gz
	mkdir -p $(RESULTS_ROOT)/$*/
	cp $^ $@		

$(RESULTS_ROOT)/%/lat.txt: $(RESULTS_ROOT)/%/lat.gz | steps utils
	tsk=prepareLat;id=$*;. $(APP_DIR)/m_start.sh $@

	gunzip -c $^ | lattice-copy ark:- ark,t:- | \
		$(APP_DIR)/int.to.word -c 3 -v $(MODELS_ROOT)/words.txt > $@

$(RESULTS_ROOT)/%/lat.nb10.txt: $(RESULTS_ROOT)/%/lat.gz | steps utils
	tsk=lattice2n10;id=$*;. $(APP_DIR)/m_start.sh $@
	gunzip -c $^ | \
	  lattice-to-nbest --n=10 ark:- ark:- | \
		lattice-copy ark:- ark,t:- | \
		$(APP_DIR)/int.to.word -c 3 -v $(MODELS_ROOT)/words.txt > $@		

#restore lattice
$(OUTPUT_ROOT)/trans/%/restore/.done: $(OUTPUT_ROOT)/trans/%/restore/lat.scaled.gz | steps utils
	tsk=restoreLat;id=$*;. $(APP_DIR)/m_start.sh $@

	unk_id=$$(grep '<unk>' $(MODELS_ROOT)/words.txt | cut -d' ' -f2); \
	sil_id=$$(grep '<eps>' $(MODELS_ROOT)/words.txt | cut -d' ' -f2); \
	lattice-prune --beam=7 "ark:gunzip -c $^|"  ark:- | \
		lattice-push ark:- ark:- | \
		lattice-align-words --silence-label=$${sil_id} --partial-word-label=$${unk_id} $(MODELS_ROOT)/word_boundary.int \
		$(MODELS_ROOT)/final.mdl ark:- ark,t:$(OUTPUT_ROOT)/trans/$*/restore/L1.lat || exit 1

	# Extract 1-best lattice
	# negalima naudoti lattice-to-nbest --n=1 ark,t:L1.lat ark,t:L2.lat
	# lattice-to-nbest ne tik papildo '-1' utt_id, bet ir prijungia tylas prie þodþiø pabaigos
	# taip desinchronizuodama L1 ir L2
	lattice-1best ark,t:$(OUTPUT_ROOT)/trans/$*/restore/L1.lat \
		ark,t:$(OUTPUT_ROOT)/trans/$*/restore/L2.lat || exit 1

	# Replace transition-ids by phone-ids (timing is lost)
	lattice-to-phone-lattice --replace-words=false $(MODELS_ROOT)/final.mdl \
		ark,t:$(OUTPUT_ROOT)/trans/$*/restore/L1.lat \
		ark,t:$(OUTPUT_ROOT)/trans/$*/restore/L3.lat  || exit 1

	# Perform processing 
	cd restore && perl lat_restore.pl $(OUTPUT_ROOT)/trans/$*/restore/L1.lat \
		$(OUTPUT_ROOT)/trans/$*/restore/L2.lat \
		$(OUTPUT_ROOT)/trans/$*/restore/L3.lat $(MODELS_ROOT)/words.txt \
		$(MODELS_ROOT)/phones.txt --join-spk --join-num 0.03 > $(OUTPUT_ROOT)/trans/$*/restore/lat.restored.txt || exit 1
	touch $@ 

$(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt: $(OUTPUT_ROOT)/trans/%/restore/.done
	tsk=punctuate;id=$*;. $(APP_DIR)/m_start.sh $@

	mkdir -p $(OUTPUT_ROOT)/trans/$*/punctuated
ifeq ($(PUNCTUATE), true)
	$(APP_DIR)/punct.lattice -u $(PR_URL) $(OUTPUT_ROOT)/trans/$*/restore/lat.restored.txt $@ || exit 1
else
	cat $(OUTPUT_ROOT)/trans/$*/restore/lat.restored.txt > $@
endif

$(RESULTS_ROOT)/%/.done: $(RESULTS_ROOT)/%/result.txt $(RESULTS_ROOT)/%/resultFinal.txt \
			$(RESULTS_ROOT)/%/lat.txt $(RESULTS_ROOT)/%/lat.nb10.txt  \
			$(RESULTS_ROOT)/%/webvtt.txt \
			$(RESULTS_ROOT)/%/lat.restored.txt \
			$(RESULTS_ROOT)/%/lat.restored.gz
	touch $@


