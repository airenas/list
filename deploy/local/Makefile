SUBDIRS=audio.prepare clean diarization dispatcher inform kafka.int manager \
	metrics punctuation rescore result result.make status transcription \
	upload volume web

clean_targets=$(patsubst %, %-clean, $(SUBDIRS))
dpush_targets=$(patsubst %, %-dpush, $(SUBDIRS))

.PHONY: clean-all

clean-all: $(clean_targets)

dpush-all: $(dpush_targets)

info:
	@echo "$(clean_targets)"

%-clean:
	$(MAKE) -C $* clean

%-dpush:
	$(MAKE) -C $* dpush