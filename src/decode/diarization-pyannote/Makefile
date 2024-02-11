# make file to split files according users
# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/decoded
INPUT_ROOT?=/audio.prepared

.SECONDARY:
.DELETE_ON_ERROR:

log=echo $$(date +'%Y-%m-%d %T')
logb=$(log): START
loge=$(log): END

PYTHONPATH?=/app
PYTHON_SCRIPTS_DIR?=/app/scripts
###############################################################################
DIARIZATION_APP?=python3 $(PYTHON_SCRIPTS_DIR)/diarization.py
FIX_SIL_APP?=python3 $(PYTHON_SCRIPTS_DIR)/fix_short_silences.py
RTTM2SEG_APP?=python3 $(PYTHON_SCRIPTS_DIR)/rttm_to_seg.py
FIX_VAD_SIL_APP?=python3 $(PYTHON_SCRIPTS_DIR)/fix_vad_silences.py
VAD_APP?=python3 $(PYTHON_SCRIPTS_DIR)/vad.py
FIX_SEGMENTS_APP?=$(APP_DIR)/fix.segments
SEND_METRIC_APP?=${APP_DIR}/send.metric
JOIN_SIL?=0.5
JOIN_GAP?=3
###############################################################################
## for metrics
SHELL=$(APP_DIR)/m_shell.sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
worker=diarization-pyannote
###############################################################################

# Speaker diarization
$(OUTPUT_ROOT)/diarization/%/res.rttm: $(INPUT_ROOT)/%/prepared.wav
	tsk=split;id=$*;. $(APP_DIR)/m_start.sh $@

	echo "Number of speakers = $(NUMBER_OF_SPEAKERS)"
	rm -rf `dirname $@`
	mkdir -p `dirname $@`
	if [ -z "$(NUMBER_OF_SPEAKERS)" ]; \
		then diarization_opts=""; \
		else diarization_opts="--num-speakers=$(NUMBER_OF_SPEAKERS)"; fi; \
	echo "diarization_opts=$$diarization_opts"; \
	$(DIARIZATION_APP) --input=$^ --output=$@ $$diarization_opts

# additional VAD
$(OUTPUT_ROOT)/diarization/%/vad.rttm: $(INPUT_ROOT)/%/prepared.wav
	tsk=vad;id=$*;. $(APP_DIR)/m_start.sh $@
	$(VAD_APP) --input=$^ --output=$@ --rate=16000

# Fix silences
$(OUTPUT_ROOT)/diarization/%/res.sil.rttm: $(OUTPUT_ROOT)/diarization/%/res.rttm 
	tsk=rttm-fix;id=$*;. $(APP_DIR)/m_start.sh $@
	$(FIX_SIL_APP) --input=$^ --output=$@ --len="$$(cat $(INPUT_ROOT)/$*/len)" --join-sil $(JOIN_SIL) --join-gap $(JOIN_GAP)
# add speech from vad
$(OUTPUT_ROOT)/diarization/%/res.fix.rttm: $(OUTPUT_ROOT)/diarization/%/res.sil.rttm $(OUTPUT_ROOT)/diarization/%/vad.rttm
	tsk=rttm-vad-fix;id=$*;. $(APP_DIR)/m_start.sh $@
	$(FIX_VAD_SIL_APP) --input=$(OUTPUT_ROOT)/diarization/$*/res.sil.rttm --vad=$(OUTPUT_ROOT)/diarization/$*/vad.rttm --output=$@ \
			--len="$$(cat $(INPUT_ROOT)/$*/len)" --join-sil $(JOIN_SIL) --join-gap $(JOIN_GAP)

# To kaldi format
$(OUTPUT_ROOT)/diarization/%/show.seg: $(OUTPUT_ROOT)/diarization/%/res.fix.rttm
	tsk=rttm2seg;id=$*;. $(APP_DIR)/m_start.sh $@
	$(RTTM2SEG_APP) --input=$^ --output=$@

$(OUTPUT_ROOT)/diarization/%/show.fixed.seg: $(OUTPUT_ROOT)/diarization/%/show.seg 
	$(FIX_SEGMENTS_APP) $< $@


$(OUTPUT_ROOT)/trans/%/wav.scp: $(OUTPUT_ROOT)/diarization/%/show.fixed.seg 
	$(logb) $@
	mkdir -p `dirname $@`
	cat $< | awk 'NR==1 { printf "prepared %s/%s/prepared.wav", "$(INPUT_ROOT)", "$*" }' | LC_ALL=C sort > $@
	$(loge) $@

$(OUTPUT_ROOT)/trans/%/utt2spk: $(OUTPUT_ROOT)/diarization/%/show.fixed.seg 
	$(logb) $@
	mkdir -p `dirname $@`
	cat $< | \
		awk '{ printf "%s-S%04d---%09.3f-%09.3f %s-S%04d\n", \
			$$1, substr($$8,2), $$3/100.0, ($$3+$$4)/100.0, $$1, substr($$8,2)}' | \
		LC_ALL=C sort > $@
	$(loge) $@		

$(OUTPUT_ROOT)/trans/%/segments: $(OUTPUT_ROOT)/diarization/%/show.fixed.seg
	$(logb) $@
	mkdir -p `dirname $@`
	cat $< | \
		awk '{ printf "%s-S%04d---%09.3f-%09.3f %s %.3f %.3f\n", \
			$$1, substr($$8,2), $$3/100.0, ($$3+$$4)/100.0, $$1, $$3/100.0, ($$3+$$4)/100.0}' | \
		LC_ALL=C sort > $@
	$(loge) $@		

$(OUTPUT_ROOT)/trans/%/spk2utt: $(OUTPUT_ROOT)/trans/%/utt2spk
	$(logb) $@
	$(APP_DIR)/utt2spk_to_spk2utt.pl $^ > $@ || exit 1
	$(loge) $@

$(OUTPUT_ROOT)/trans/%/.diarization.done: $(OUTPUT_ROOT)/trans/%/spk2utt $(OUTPUT_ROOT)/trans/%/segments \
											$(OUTPUT_ROOT)/trans/%/wav.scp
	touch $@


info:
	echo OUTPUT_ROOT=$(OUTPUT_ROOT)
	echo PYTHONPATH=$(PYTHONPATH)
	echo PYTHON_SCRIPTS_DIR=$(PYTHON_SCRIPTS_DIR)
	echo DIARIZATION_APP=$(DIARIZATION_APP)
	echo APP_DIR=$(APP_DIR)
	echo SHELL=$(SHELL)
	