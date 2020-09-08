include Makefile.options.local

SUBDIRS=audio.prepare clean diarization diarization-alt dispatcher inform kafka.int manager \
	metrics punctuation rescore result result.make status transcription \
	upload volume web

clean_targets=$(patsubst %, %-clean, $(SUBDIRS))
dpush_targets=$(patsubst %, %-dpush, $(SUBDIRS))

.PHONY: clean-all

clean-all: $(clean_targets)

dpush-all: $(dpush_targets)

info:
	@echo "$(clean_targets)"

services:
	curl -X GET https://$(semantika_cred)@semantikadocker.vdu.lt/v2/_catalog | jq

versions:
	@$(if $(strip $(srv)),echo "Service = $(srv)",echo No srv set && exit 1)
	curl -X GET https://$(semantika_cred)@semantikadocker.vdu.lt/v2/trans/$(srv)/tags/list | jq

%-clean:
	$(MAKE) -C $* clean

%-dpush:
	$(MAKE) -C $* dpush