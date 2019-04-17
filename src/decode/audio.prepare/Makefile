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

$(OUTPUT_ROOT):
	mkdir $(OUTPUT_ROOT)

$(OUTPUT_ROOT)/%.wav: $(INPUT_ROOT)/%.wav $(OUTPUT_ROOT)
	sox $< -c 1 -b 16 $@ rate -v 16k

$(OUTPUT_ROOT)/%.wav: $(INPUT_ROOT)/%.mp3 $(OUTPUT_ROOT)
	ffmpeg -i $< -f sox - | sox -t sox - -c 1 -b 16 $@ rate -v 16k	

$(OUTPUT_ROOT)/%.wav: $(INPUT_ROOT)/%.mp4 $(OUTPUT_ROOT)
	ffmpeg -i $< -f sox - | sox -t sox - -c 1 -b 16 $@ rate -v 16k

test: $(OUTPUT_ROOT)/a.wav