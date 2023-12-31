# make file to convert audio to wav 16k from various formats
# convert from 
# - wav
# - mp3
# - mp4

# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/audio.prepared
INPUT_ROOT?=/audio.in
###############################################################################
DENOISE_APP?=${APP_DIR}/sfg
SEND_METRIC_APP?=${APP_DIR}/send.metric
###############################################################################
# in seconds
MIN_AUDIO_LEN?=0.5
MAX_AUDIO_LEN?=7200
###############################################################################
CHECK_FORMAT_CMD=r=$$?;\
		if [ $$r -gt 0 ]; then printf "[[[%s:WRONG_FORMAT]]]\n" ErrorCode >&2; exit 2; fi
CHECK_LEN_CMD=ok=$$(echo "$(MIN_AUDIO_LEN) > $$(sox --i -D $@)" | bc -l) ;\
		if [ $$ok -gt 0 ] ; then printf "[[[%s:TOO_SHORT_AUDIO]]]\n" ErrorCode >&2; exit 3; fi;\
		ok=$$(echo "$(MAX_AUDIO_LEN) < $$(sox --i -D $@)" | bc -l) ;\
		if [ $$ok -gt 0 ] ; then printf "[[[%s:TOO_LONG_AUDIO]]]\n" ErrorCode >&2; exit 3; fi
###############################################################################
## for metrics
SHELL=$(APP_DIR)/m_shell.sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
.SECONDARY:
worker=audio.prepare
###############################################################################

$(OUTPUT_ROOT):
	mkdir -p $(OUTPUT_ROOT)

$(OUTPUT_ROOT)/%/converted.wav: $(INPUT_ROOT)/%.wav | $(OUTPUT_ROOT)
	tsk=wav;id=$*;. $(APP_DIR)/m_start.sh $@
	mkdir -p `dirname $@`
	sox $< -c 1 -b 16 $@ rate -v 16k 
	$(CHECK_FORMAT_CMD)
	$(CHECK_LEN_CMD)

$(OUTPUT_ROOT)/%/converted.wav: $(INPUT_ROOT)/%.mp3 | $(OUTPUT_ROOT)
	$(MAKE) $@ a_ext=mp3
$(OUTPUT_ROOT)/%/converted.wav: $(INPUT_ROOT)/%.mp4 | $(OUTPUT_ROOT)
	$(MAKE) $@ a_ext=mp4
$(OUTPUT_ROOT)/%/converted.wav: $(INPUT_ROOT)/%.m4a | $(OUTPUT_ROOT)
	$(MAKE) $@ a_ext=m4a
$(OUTPUT_ROOT)/%/converted.wav: $(INPUT_ROOT)/%.ogg | $(OUTPUT_ROOT)
	$(MAKE) $@ a_ext=ogg	
$(OUTPUT_ROOT)/%/converted.wav: $(INPUT_ROOT)/%.wma | $(OUTPUT_ROOT)
	$(MAKE) $@ a_ext=wma	
$(OUTPUT_ROOT)/%/converted.wav: $(INPUT_ROOT)/%.webm | $(OUTPUT_ROOT)
	$(MAKE) $@ a_ext=webm		

$(OUTPUT_ROOT)/%/converted.wav: $(INPUT_ROOT)/%.$(a_ext) | $(OUTPUT_ROOT)
	tsk=$(a_ext);id=$*;. $(APP_DIR)/m_start.sh $@
	mkdir -p `dirname $@`
	ffmpeg -i $< -f sox - | sox -t sox - -c 1 -b 16 $@ rate -v 16k 
	$(CHECK_FORMAT_CMD)
	$(CHECK_LEN_CMD)

###############################################################################
$(OUTPUT_ROOT)/%/prepared.wav: $(OUTPUT_ROOT)/%/converted.wav
ifeq ($(DENOISE), true)
	tsk=denoise;id=$*;. $(APP_DIR)/m_start.sh $@
	echo "main ()" > $(OUTPUT_ROOT)/$*/denoise.cfg
	echo "{" >> $(OUTPUT_ROOT)/$*/denoise.cfg
	echo '	olasynth("$^", 0, "$@");' >> $(OUTPUT_ROOT)/$*/denoise.cfg
	echo "}" >> $(OUTPUT_ROOT)/$*/denoise.cfg
	echo "Run denoise"
	${DENOISE_APP} $(OUTPUT_ROOT)/$*/denoise.cfg -nofinish
else
	cp $^ $@
endif
###############################################################################
$(OUTPUT_ROOT)/%/len: $(OUTPUT_ROOT)/%/prepared.wav
	sox --i -D $^ > $@
	echo Audio len: $$(cat $@)
###############################################################################
$(OUTPUT_ROOT)/%/.done: $(OUTPUT_ROOT)/%/prepared.wav $(OUTPUT_ROOT)/%/len
	touch $@
###############################################################################
