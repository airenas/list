# make file to split files according users
# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/decoded
INPUT_ROOT?=/audio.prepared
KALDI_ROOT?=/kaldi

.SECONDARY:
.DELETE_ON_ERROR:

log=echo $$(date +'%Y-%m-%d %T')
logb=$(log): START
loge=$(log): END

###############################################################################
## for metrics
SHELL=$(APP_DIR)/m_shell.sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
worker=diarization
###############################################################################

# Speaker diarization
$(OUTPUT_ROOT)/diarization/%/show.seg: $(INPUT_ROOT)/%/prepared.wav
	tsk=split;id=$*;. $(APP_DIR)/m_start.sh $@

	echo "Number of speakers = $(NUMBER_OF_SPEAKERS)"
	rm -rf `dirname $@`
	mkdir -p `dirname $@`
	if [ -z "$(NUMBER_OF_SPEAKERS)" ]; \
		then diarization_opts="--min_speakers=1 --max_speakers=100"; \
		else diarization_opts="--min_speakers=$(NUMBER_OF_SPEAKERS) --max_speakers=$(NUMBER_OF_SPEAKERS)"; fi; \
	echo "diarization_opts=$$diarization_opts"; \
	(cd /app; python3 /app/main.py --wav_path=$^ --seg_path=$@ $$diarization_opts)

$(OUTPUT_ROOT)/diarization/%/show.fixed.seg: $(OUTPUT_ROOT)/diarization/%/show.seg 
	$(APP_DIR)/fix.segments $< $@


$(OUTPUT_ROOT)/trans/%/wav.scp: $(OUTPUT_ROOT)/diarization/%/show.fixed.seg 
	$(logb) $@
	mkdir -p `dirname $@`
	cat $< | awk 'NR==1 { printf "%s %s/%s/prepared.wav", $$1, "$(INPUT_ROOT)", $$1 }' | LC_ALL=C sort > $@
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
	$(KALDI_ROOT)/utils/utt2spk_to_spk2utt.pl $^ > $@
	$(loge) $@

$(OUTPUT_ROOT)/trans/%/.diarization.done: $(OUTPUT_ROOT)/trans/%/spk2utt $(OUTPUT_ROOT)/trans/%/segments \
											$(OUTPUT_ROOT)/trans/%/wav.scp
	touch $@
