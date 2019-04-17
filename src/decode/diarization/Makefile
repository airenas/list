# make file to split files according users
# copied and modified from https://github.com/alumae/kaldi-offline-transcriber
# see Licenses/LICENSE.alumae

-include Makefile.options

OUTPUT_ROOT?=/decoded
INPUT_ROOT?=/audio.prepared
APPS_ROOT?=/diarization
LOCALCLASSPATH?=$(APPS_ROOT)/lib/LIUM_SpkDiarization-4.2.jar
DO_MUSIC_DETECTION=no
JAVA_BIN?=java

.SECONDARY:
.DELETE_ON_ERROR:

# Speaker diarization
$(OUTPUT_ROOT)/diarization/%/show.seg: $(INPUT_ROOT)/%.wav
	rm -rf `dirname $@`
	mkdir -p `dirname $@`
	echo "$* 1 0 1000000000 U U U 1" >  `dirname $@`/show.uem.seg;
	if [ $(DO_MUSIC_DETECTION) = yes ]; then diarization_opts="-m"; fi; \
	(export LOCALCLASSPATH=$(LOCALCLASSPATH); export JAVA_BIN=$(JAVA_BIN); export MODELS_ROOT=$(MODELS_ROOT); $(APPS_ROOT)/diarization.sh $$diarization_opts $^ `dirname $@`/show.uem.seg)

$(OUTPUT_ROOT)/audio/segmented/%: $(OUTPUT_ROOT)/diarization/%/show.seg
	rm -rf $@
	mkdir -p $@
	cat $^ | cut -f 3,4,8 -d " " | \
	while read LINE ; do \
		start=`echo $$LINE | cut -f 1 -d " " | perl -npe '$$_=$$_/100.0'`; \
		len=`echo $$LINE | cut -f 2 -d " " | perl -npe '$$_=$$_/100.0'`; \
		sp_id=`echo $$LINE | cut -f 3 -d " "`; \
		timeformatted=`echo "$$start $$len" | perl -ne '@t=split(); $$start=$$t[0]; $$len=$$t[1]; $$end=$$start+$$len; printf("%08.3f-%08.3f\n", $$start,$$end);'` ; \
		sox $(INPUT_ROOT)/$*.wav --norm $@/$*_$${timeformatted}_$${sp_id}.wav trim $$start $$len ; \
	done

$(OUTPUT_ROOT)/trans/%/wav.scp: $(OUTPUT_ROOT)/audio/segmented/%
	mkdir -p `dirname $@`
	/bin/ls $</*.wav  | \
		perl -npe 'chomp; $$orig=$$_; s/.*\/(.*)_(\d+\.\d+-\d+\.\d+)_(S\d+)\.wav/\1-\3---\2/; $$_=$$_ .  " $$orig\n";' | LC_ALL=C sort > $@

$(OUTPUT_ROOT)/trans/%/utt2spk: $(OUTPUT_ROOT)/trans/%/wav.scp
	cat $^ | perl -npe 's/\s+.*//; s/((.*)---.*)/\1 \2/' > $@

$(OUTPUT_ROOT)/trans/%/spk2utt: $(OUTPUT_ROOT)/trans/%/utt2spk
	$(APPS_ROOT)/utt2spk_to_spk2utt.pl $^ > $@

test: $(OUTPUT_ROOT)/trans/a/spk2utt
