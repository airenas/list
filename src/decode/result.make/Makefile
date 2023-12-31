# make file to prepare the result from lattice
# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/decoded
RESULTS_ROOT?=/results
APPS_ROOT?=/kaldi
FINAL_PASS=2-pass
APP_DIR?=/app

###############################################################################
LATTICE_TO_TEXT_APP?=$(APP_DIR)/lattice.to.text
EMPTY_LATTICE_APP?=$(APP_DIR)/empty.lattice
FIX_LATTICE_TIME_APP?=$(APP_DIR)/fix.lattice.time
INT_TO_WORD_APP?=$(APP_DIR)/int.to.word
LATTICE_TO_WEBVTT_APP?=$(APP_DIR)/lattice.to.webvtt
SEND_METRIC_APP?=${APP_DIR}/send.metric
###############################################################################

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

###############################################################################
#scale
$(OUTPUT_ROOT)/trans/%/restore/lat.scaled.gz: $(OUTPUT_ROOT)/trans/%/$(FINAL_PASS)/.done
	tsk=scale;id=$*;. $(APP_DIR)/m_start.sh $@

	mkdir -p $(OUTPUT_ROOT)/trans/$*/restore
ifeq ($(RESULT_LMWT),)
	opt_lmwt=$$(cat "$(MODELS_ROOT)/lmwt")
else
	opt_lmwt=$(RESULT_LMWT)
	echo "Use opt_lmwt from env settings"
endif
ifeq ($(RESULT_WIP),)
	opt_wip=$$(cat "$(MODELS_ROOT)/wip")
else
	opt_wip=$(RESULT_WIP)
	echo "Use opt_wip from env settings"
endif
	echo "opt_lmwt=$${opt_lmwt}"
	echo "opt_wip=$${opt_wip}"
	lattice-scale --inv-acoustic-scale=$${opt_lmwt} \
		"ark:gunzip -c $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode/lat.1.gz|"  ark:- | \
		lattice-add-penalty --word-ins-penalty=$${opt_wip} ark:- "ark:| gzip -c > $@"

###############################################################################
ifeq ($(NO_SPEECH), true)
$(RESULTS_ROOT)/%/result.txt: 
	mkdir -p $(RESULTS_ROOT)/$*/
	touch $@
else
$(RESULTS_ROOT)/%/result.txt: $(OUTPUT_ROOT)/trans/%/restore/.done
	tsk=lattice2text;id=$*;. $(APP_DIR)/m_start.sh $@
	
	mkdir -p $(RESULTS_ROOT)/$*/
	$(LATTICE_TO_TEXT_APP) $(OUTPUT_ROOT)/trans/$*/restore/lat.restored.txt $@
endif	
###############################################################################

$(RESULTS_ROOT)/%/resultFinal.txt: $(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt
	tsk=lattice2textFinal;id=$*;. $(APP_DIR)/m_start.sh $@

	mkdir -p $(RESULTS_ROOT)/$*/
	$(LATTICE_TO_TEXT_APP) $^ $@

$(RESULTS_ROOT)/%/webvtt.txt: $(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt
	tsk=lattice2webvtt;id=$*;. $(APP_DIR)/m_start.sh $@
	
	mkdir -p $(RESULTS_ROOT)/$*/
	$(LATTICE_TO_WEBVTT_APP) $^ $@

###############################################################################
$(RESULTS_ROOT)/%/lat.restored.txt: $(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt
	mkdir -p $(RESULTS_ROOT)/$*/
ifeq ($(NO_SPEECH), true)
	$(EMPTY_LATTICE_APP) -l $$(cat $(OUTPUT_ROOT)/audio/$*/len) -s '[]' $@
else
	$(FIX_LATTICE_TIME_APP) -l $$(cat $(OUTPUT_ROOT)/audio/$*/len) -s '[]' $^ $@
endif
###############################################################################

$(RESULTS_ROOT)/%/lat.restored.gz: $(RESULTS_ROOT)/%/lat.restored.txt
	mkdir -p $(RESULTS_ROOT)/$*/
	gzip -c $^ > $@

###############################################################################
# lat.gz
ifeq ($(NO_SPEECH), true)
$(RESULTS_ROOT)/%/lat.gz: 
	mkdir -p $(RESULTS_ROOT)/$*/
	touch $@
else
$(RESULTS_ROOT)/%/lat.gz: $(OUTPUT_ROOT)/trans/%/restore/lat.scaled.gz
	mkdir -p $(RESULTS_ROOT)/$*/
	cp $^ $@		
endif	
###############################################################################

$(RESULTS_ROOT)/%/lat.txt: $(RESULTS_ROOT)/%/lat.gz | steps utils
	tsk=prepareLat;id=$*;. $(APP_DIR)/m_start.sh $@

	gunzip -c $^ | lattice-copy ark:- ark,t:- | \
		$(INT_TO_WORD_APP) -f 3 -v $(MODELS_ROOT)/words.txt > $@

$(RESULTS_ROOT)/%/lat.nb10.txt: $(RESULTS_ROOT)/%/lat.gz | steps utils
	tsk=lattice2n10;id=$*;. $(APP_DIR)/m_start.sh $@
	gunzip -c $^ | \
	  lattice-to-nbest --n=10 ark:- ark:- | \
		lattice-copy ark:- ark,t:- | \
		$(INT_TO_WORD_APP) -f 3 -v $(MODELS_ROOT)/words.txt > $@		

#restore lattice
$(OUTPUT_ROOT)/trans/%/restore/.done: $(OUTPUT_ROOT)/trans/%/restore/lat.scaled.gz | steps utils
	tsk=restoreLat;id=$*;. $(APP_DIR)/m_start.sh $@
	
	echo Scripts dir=$(SCRIPTS_DIR)
ifneq ("$(wildcard $(SCRIPTS_DIR)/restore.sh)","")
	echo Running: $(SCRIPTS_DIR)/restore.sh
	$(SCRIPTS_DIR)/restore.sh $(MODELS_ROOT) $^ $(OUTPUT_ROOT)/trans/$*/restore || exit 1
else
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

	numJoinParam='--join-num 0.03'; if [ "$(SKIP_NUM_JOIN)" == "1" ]; then numJoinParam=''; fi
	spkJoinParam='--join-spk'; if [ "$(SKIP_SPEAKER_JOIN)" == "true" ] || [ "$(SKIP_SPEAKER_JOIN)" == "1" ]; then spkJoinParam=''; fi
	echo "SKIP_NUM_JOIN=$(SKIP_NUM_JOIN); numJoinParam=$${numJoinParam}"
	echo "SKIP_SPEAKER_JOIN=$(SKIP_SPEAKER_JOIN); spkJoinParam=$${spkJoinParam}"
	# Perform processing 
	cd restore && perl lat_restore.pl $(OUTPUT_ROOT)/trans/$*/restore/L1.lat \
		$(OUTPUT_ROOT)/trans/$*/restore/L2.lat \
		$(OUTPUT_ROOT)/trans/$*/restore/L3.lat $(MODELS_ROOT)/words.txt \
		$(MODELS_ROOT)/phones.txt $${spkJoinParam} $${numJoinParam} > $(OUTPUT_ROOT)/trans/$*/restore/lat.restored.txt || exit 1
endif
	touch $@ 

###############################################################################
# punctuate task
###############################################################################
ifeq ($(NO_SPEECH), true)
$(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt:
	echo Generating fake lat.restored.txt
	mkdir -p $(OUTPUT_ROOT)/trans/$*/punctuated
	$(EMPTY_LATTICE_APP) -l $$(cat $(OUTPUT_ROOT)/audio/$*/len) -s '<eps>' $@
else
###############################################################################
$(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt: $(OUTPUT_ROOT)/trans/%/restore/.done
	tsk=punctuate;id=$*;. $(APP_DIR)/m_start.sh $@

	mkdir -p $(OUTPUT_ROOT)/trans/$*/punctuated
###############################################################################
ifeq ($(PUNCTUATE), true)
	$(APP_DIR)/punct.lattice -u $(PR_URL) $(OUTPUT_ROOT)/trans/$*/restore/lat.restored.txt $@ || exit 1
else
	cat $(OUTPUT_ROOT)/trans/$*/restore/lat.restored.txt > $@
endif
###############################################################################
endif
###############################################################################

$(RESULTS_ROOT)/%/.done: $(RESULTS_ROOT)/%/result.txt $(RESULTS_ROOT)/%/resultFinal.txt \
			$(RESULTS_ROOT)/%/lat.txt $(RESULTS_ROOT)/%/lat.nb10.txt  \
			$(RESULTS_ROOT)/%/webvtt.txt \
			$(RESULTS_ROOT)/%/lat.restored.txt \
			$(RESULTS_ROOT)/%/lat.restored.gz
	touch $@
