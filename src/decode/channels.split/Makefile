# make file to split audio file according left and right channel
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
worker=channels.split
###############################################################################
$(OUTPUT_ROOT)/$(id):
	mkdir -p $@	
###############################################################################
$(OUTPUT_ROOT)/$(id)/audio_%_channel.wav: $(INPUT_ROOT)/$(id).wav | $(OUTPUT_ROOT)/$(id)
	tsk=to=split;id=$(id);. $(APP_DIR)/m_start.sh $@
	ffmpeg -i $<  -filter:a 'pan=mono|FC=$*' $@
$(OUTPUT_ROOT)/$(id)/audio_%_channel.wav: $(INPUT_ROOT)/$(id).mp3 | $(OUTPUT_ROOT)/$(id)
	tsk=to=split;id=$(id);. $(APP_DIR)/m_start.sh $@
	ffmpeg -i $<  -filter:a 'pan=mono|FC=$*' $@	
$(OUTPUT_ROOT)/$(id)/audio_%_channel.wav: $(INPUT_ROOT)/$(id).mp4 | $(OUTPUT_ROOT)/$(id)
	tsk=to=split;id=$(id);. $(APP_DIR)/m_start.sh $@
	ffmpeg -i $<  -filter:a 'pan=mono|FC=$*' $@
$(OUTPUT_ROOT)/$(id)/audio_%_channel.wav: $(INPUT_ROOT)/$(id).m4a | $(OUTPUT_ROOT)/$(id)
	tsk=to=split;id=$(id);. $(APP_DIR)/m_start.sh $@
	ffmpeg -i $<  -filter:a 'pan=mono|FC=$*' $@	
###############################################################################
$(OUTPUT_ROOT)/$(id)/.split.done: $(OUTPUT_ROOT)/$(id)/audio_FL_channel.wav \
	$(OUTPUT_ROOT)/$(id)/audio_FR_channel.wav
	touch $@
###############################################################################
