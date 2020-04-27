# make file to split files according users
# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/decoded
INPUT_ROOT?=/audio.prepared
APPS_ROOT?=/diarization
KALDI_ROOT?=/kaldi
LOCALCLASSPATH?=$(APPS_ROOT)/lib/LIUM_SpkDiarization-4.2.jar
DO_MUSIC_DETECTION=no
JAVA_BIN?=java

.SECONDARY:
.DELETE_ON_ERROR:

log=@echo $$(date +'%Y-%m-%d %T')
logb=$(log): START
loge=$(log): END

# Speaker diarization
$(OUTPUT_ROOT)/diarization/%/show.seg: $(INPUT_ROOT)/%.wav
	$(logb) $@
	@echo "Number of speakers = $(NUMBER_OF_SPEAKERS)"
	rm -rf `dirname $@`
	mkdir -p `dirname $@`
	echo "$* 1 0 1000000000 U U U 1" >  `dirname $@`/show.uem.seg;
	if [ $(DO_MUSIC_DETECTION) = yes ]; then diarization_opts="-m"; fi; \
	(export LOCALCLASSPATH=$(LOCALCLASSPATH); export JAVA_BIN=$(JAVA_BIN); export MODELS_ROOT=$(MODELS_ROOT); $(APPS_ROOT)/diarization.sh $$diarization_opts $^ `dirname $@`/show.uem.seg)
	$(loge) $@

$(OUTPUT_ROOT)/trans/%/wav.scp: $(OUTPUT_ROOT)/diarization/%/show.seg 
	$(logb) $@
	@mkdir -p `dirname $@`
	cat $@ | awk 'NR==1 { printf "%s %s/%s.wav", $$1, "$(INPUT_ROOT)", $$1 }' | LC_ALL=C sort > $@
	$(loge) $@

$(OUTPUT_ROOT)/trans/%/utt2spk: $(OUTPUT_ROOT)/diarization/%/show.seg 
	$(logb) $@
	@mkdir -p `dirname $@`
	cat $@ | \
		awk '{ printf "%s-S%04d---%09.3f-%09.3f %s-S%04d\n", \
			$$1, substr($$8,2), $$3/100.0, ($$3+$$4)/100.0, $$1, substr($$8,2)}' | \
		LC_ALL=C sort > $@
	$(loge) $@		

$(OUTPUT_ROOT)/trans/%/segments: $(OUTPUT_ROOT)/diarization/%/show.seg
	$(logb) $@
	@mkdir -p `dirname $@`
	cat $@ | \
		awk '{ printf "%s-S%04d---%09.3f-%09.3f %s %.3f %.3f\n", \
			$$1, substr($$8,2), $$3/100.0, ($$3+$$4)/100.0, $$1, $$3/100.0, ($$3+$$4)/100.0}' | \
		LC_ALL=C sort > $@
	$(loge) $@		

$(OUTPUT_ROOT)/trans/%/spk2utt: $(OUTPUT_ROOT)/trans/%/utt2spk
	$(logb) $@
	$(KALDI_ROOT)/utils/utt2spk_to_spk2utt.pl $^ > $@
	$(loge) $@

$(OUTPUT_ROOT)/trans/%/.diarization.done: $(OUTPUT_ROOT)/trans/%/utt2spk $(OUTPUT_ROOT)/trans/%/segments \
											$(OUTPUT_ROOT)/trans/%/wav.scp
	touch $@

