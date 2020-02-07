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

WRONG_FORMAT_CMD=(echo -n [[[ >&2 && echo ErrorCode:WRONG_FORMAT]]] >&2 && exit 2)
TOO_SHORT_AUDIO_CMD=(echo -n [[[ >&2 && echo ErrorCode:TOO_SHORT_AUDIO]]] >&2 && exit 3)

define checkAudioLen
	{ \
		ok=$$(echo "$(MIN_LEN) > $$(sox --i -D $(1))" | bc -l) ;\
		if [ $$ok -gt 0 ] ;\
		then \
			$(TOO_SHORT_AUDIO_CMD);\
		fi ;\
    }
endef

log=@echo $$(date +'%Y-%m-%d %T')
logb=$(log): START
loge=$(log): END

$(OUTPUT_ROOT):
	mkdir $(OUTPUT_ROOT)

$(OUTPUT_ROOT)/%.wav: $(INPUT_ROOT)/%.wav $(OUTPUT_ROOT)
	$(logb) $@
	sox $< -c 1 -b 16 $@ rate -v 16k || $(WRONG_FORMAT_CMD)
	$(call checkAudioLen,$@)
	$(loge) $@

$(OUTPUT_ROOT)/%.wav: $(INPUT_ROOT)/%.mp3 $(OUTPUT_ROOT)
	$(logb) $@
	ffmpeg -i $< -f sox - | sox -t sox - -c 1 -b 16 $@ rate -v 16k || $(WRONG_FORMAT_CMD)
	$(call checkAudioLen,$@)
	$(loge) $@

$(OUTPUT_ROOT)/%.wav: $(INPUT_ROOT)/%.mp4 $(OUTPUT_ROOT)
	$(logb) $@
	ffmpeg -i $< -f sox - | sox -t sox - -c 1 -b 16 $@ rate -v 16k || $(WRONG_FORMAT_CMD)
	$(call checkAudioLen,$@)
	$(loge) $@

test: $(OUTPUT_ROOT)/a.wav
