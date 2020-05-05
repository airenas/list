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


MIN_LEN?=0.5

CHECK_FORMAT_CMD=r=$$?;\
		if [ $$r -gt 0 ]; then printf "[[[%s:WRONG_FORMAT]]]\n" ErrorCode >&2; exit 2; fi
CHECK_LEN_CMD=ok=$$(echo "$(MIN_LEN) > $$(sox --i -D $@)" | bc -l) ;\
		if [ $$ok -gt 0 ] ; then printf "[[[%s:TOO_SHORT_AUDIO]]]\n" ErrorCode >&2; exit 3; fi

###############################################################################
## for metrics
SHELL=$(SCRIPTS_DIR)/m_shell.sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
worker=audio.prepare
###############################################################################

$(OUTPUT_ROOT):
	mkdir $(OUTPUT_ROOT)

$(OUTPUT_ROOT)/%.wav: $(INPUT_ROOT)/%.wav | $(OUTPUT_ROOT)
	tsk=wav;id=$*;. $(SCRIPTS_DIR)/m_start.sh $@
	sox $< -c 1 -b 16 $@ rate -v 16k 
	$(CHECK_FORMAT_CMD)
	$(CHECK_LEN_CMD)

$(OUTPUT_ROOT)/%.wav: $(INPUT_ROOT)/%.mp3 | $(OUTPUT_ROOT)
	tsk=mp3;id=$*;. $(SCRIPTS_DIR)/m_start.sh $@
	
	ffmpeg -i $< -f sox - | sox -t sox - -c 1 -b 16 $@ rate -v 16k
	$(CHECK_FORMAT_CMD)
	$(CHECK_LEN_CMD)

$(OUTPUT_ROOT)/%.wav: $(INPUT_ROOT)/%.mp4 | $(OUTPUT_ROOT)
	tsk=mp4;id=$*;. $(SCRIPTS_DIR)/m_start.sh $@
	ffmpeg -i $< -f sox - | sox -t sox - -c 1 -b 16 $@ rate -v 16k 
	$(CHECK_FORMAT_CMD)
	$(CHECK_LEN_CMD)
