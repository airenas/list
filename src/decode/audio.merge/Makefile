# make file to merge audio streams to one mp3 file
-include Makefile.options

INPUT_ROOT?=/audio.in
OUTPUT_ROOT?=/audio.in
###############################################################################
SEND_METRIC_APP?=${APP_DIR}/send.metric
###############################################################################
## for metrics
SHELL=$(APP_DIR)/m_shell.sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
.SECONDARY:
worker=audio.merge
###############################################################################
files:=$(wildcard $(INPUT_ROOT)/$(id)/*.mp4 $(INPUT_ROOT)/$(id)/*.wav $(INPUT_ROOT)/$(id)/*.m4a $(INPUT_ROOT)/$(id)/*.mp3)
w_files=$(files:$(INPUT_ROOT)/$(id)/%=$(INPUT_ROOT)/$(id)/merge/%)
wav_files= $(patsubst %.mp4,%.wav,$(patsubst %.mp3,%.wav,$(patsubst %.m4a,%.wav,$(w_files))))
###############################################################################
$(INPUT_ROOT)/$(id)/merge:
	mkdir -p $@	
$(INPUT_ROOT)/$(id)/merge/%.wav: $(INPUT_ROOT)/$(id)/%.mp3 | $(INPUT_ROOT)/$(id)/merge
	ffmpeg -i $< $@
$(INPUT_ROOT)/$(id)/merge/%.wav: $(INPUT_ROOT)/$(id)/%.mp4 | $(INPUT_ROOT)/$(id)/merge
	ffmpeg -i $< $@	
$(INPUT_ROOT)/$(id)/merge/%.wav: $(INPUT_ROOT)/$(id)/%.m4a | $(INPUT_ROOT)/$(id)/merge
	ffmpeg -i $< $@		
$(INPUT_ROOT)/$(id)/merge/%.wav: $(INPUT_ROOT)/$(id)/%.wav | $(INPUT_ROOT)/$(id)/merge
	cp $< $@
###############################################################################
$(INPUT_ROOT)/$(id)/merge/$(id).wav: $(wav_files)
	tsk=to=sox;id=$(id);. $(APP_DIR)/m_start.sh $@
	sox -m $(wav_files) $@
###############################################################################
$(OUTPUT_ROOT)/$(id).mp3: $(INPUT_ROOT)/$(id)/merge/$(id).wav
	tsk=to=mp3;id=$(id);. $(APP_DIR)/m_start.sh $@
	ffmpeg -i $< $@
###############################################################################
$(OUTPUT_ROOT)/$(id)/.merge.done: $(OUTPUT_ROOT)/$(id).mp3
	touch $@
###############################################################################
