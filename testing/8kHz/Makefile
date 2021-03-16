-include Makefile.options
-include $(config)
wordfix_dir?=../../lm-scripts/wordfix
scripts_dir=../scripts
#####################################################################################
file?=1.wav
upload: 
	curl -X POST -k $(UPLOAD_SERVICE)/upload -H 'Accept: application/json' -H 'content-type: multipart/form-data' \
	-F recognizer=$(model) -F numberOfSpeakers=2 \
	-F file=@$(file)

record:
	arecord -f cd > $(file)
#####################################################################################
$(test_dir):
	mkdir -p $@
$(test_dir)/wav: 
	mkdir -p $@	
$(test_dir)/.wav.done: | $(test_dir)/wav
	ls  $(corpus_dir)/WAV | xargs -I {} sh -c "ls  $(corpus_dir)/WAV/{}/ | head -n 3 | \
		xargs -I [] cp  $(corpus_dir)/WAV/{}/[] $(test_dir)/wav/{}_[]"
	touch $@
$(test_dir)/kaldi_rec.all.txt: | $(test_dir)
	cat $(corpus_dir)/results_1/kaldiout.txt > $@
$(test_dir)/ref.all.txt: | $(test_dir)
	cat $(corpus_dir)/kaldi_data_1/textw > $@
$(test_dir)/kaldi_rec.txt: | $(test_dir)
	cat $(corpus_dir)/results_1/kaldiout.txt | grep -E "_00 |_01 |_02" > $@
$(test_dir)/ref.txt: | $(test_dir)
	cat $(corpus_dir)/kaldi_data_1/textw | grep -E "_00 |_01 |_02" > $@
$(test_dir)/%.f.txt: $(test_dir)/%.txt
	cat $< | perl $(wordfix_dir)/filter.pl $(wordfix_dir) > $@
$(test_dir)/.data.done: $(test_dir)/ref.f.txt $(test_dir)/kaldi_rec.f.txt $(test_dir)/.wav.done | $(test_dir)
	touch $@
prepare: $(test_dir)/.data.done
#####################################################################################
kaldi-error: prepare
	$(scripts_dir)/compute-wer --mode=present \
		ark:$(test_dir)/ref.f.txt  ark:$(test_dir)/kaldi_rec.f.txt
kaldi-error-all: $(test_dir)/ref.all.f.txt $(test_dir)/kaldi_rec.all.f.txt
	$(scripts_dir)/compute-wer --mode=present \
		ark:$(test_dir)/ref.all.f.txt  ark:$(test_dir)/kaldi_rec.all.f.txt		
#####################################################################################
$(test_dir)/.upload.done: $(test_dir)/.data.done
	@echo "Sending data to server $(UPLOAD_SERVICE)"
	export UPLOAD_SERVICE=$(UPLOAD_SERVICE); \
	export SKIP_NUM_JOIN=1; \
	export NUMBER_OF_SPEAKERS=$(NUMBER_OF_SPEAKERS); \
		ls -1 $(test_dir)/wav/*.wav | xargs -n1 -P10 $(scripts_dir)/send.sh $(model) > $(test_dir)/file.list; \
	touch $(test_dir)/.upload.done
	
upload-to-server: $(test_dir)/.upload.done
recognizers:
	curl -k $(UPLOAD_SERVICE)/recognizers | jq

#####################################################################################
transcription-status: $(test_dir)/file.list
	export STATUS_SERVICE=$(STATUS_SERVICE); cat $(test_dir)/file.list | xargs -n2 -P20 $(scripts_dir)/status.sh | sort	
#####################################################################################
$(test_dir)/result:
	mkdir -p $@

$(test_dir)/result/.done: $(test_dir)/file.list | $(test_dir)/result
	export RESULT_SERVICE=$(RESULT_SERVICE); cat $(test_dir)/file.list | xargs -n2 -P20 $(scripts_dir)/result.sh $(test_dir)/result
	touch $(test_dir)/result/.done
get_result: $(test_dir)/result/.done	
#####################################################################################
$(test_dir)/result/fixed:
	mkdir -p $@

$(test_dir)/recognized.txt: $(test_dir)/result/.done | $(test_dir)/result/fixed
	ls -1 $(test_dir)/result/*$(tf).txt | xargs -n 1 -I {} sh -c "f={};\
		cat \$$f | tr '\n' ' ' | sed -r 's/ +/ /g' | sed -r 's/^ //' > $(test_dir)/result/fixed/\$${f##*/}"
	awk '{print FILENAME,$$0}' $(test_dir)/result/fixed/*.txt | \
		sed 's|$(test_dir)/result/fixed/||' | sed 's|.txt||' > $@_
	mv $@_ $@	

$(test_dir)/recognized.f.txt: $(test_dir)/recognized.txt
		cat $(test_dir)/recognized.txt | perl $(wordfix_dir)/filter.pl $(wordfix_dir) > $(test_dir)/recognized.f.txt

calc-wer-1: $(test_dir)/recognized.f.txt $(test_dir)/ref.f1.txt
	$(scripts_dir)//compute-wer --mode=present \
		ark:$(test_dir)/ref.f1.txt  ark:$(test_dir)/recognized.f.txt
$(test_dir)/wer.txt: $(test_dir)/recognized.f.txt $(test_dir)/ref.f.txt
	$(scripts_dir)/compute-wer --mode=present \
		ark:$(test_dir)/ref.f.txt ark:$(test_dir)/recognized.f.txt > $(test_dir)/wer.txt 
calc-wer: $(test_dir)/wer.txt
	cat $(test_dir)/wer.txt
compare: $(test_dir)/recognized.f.txt $(test_dir)/ref.f.txt
	vimdiff $(test_dir)/ref.f.txt $(test_dir)/recognized.f.txt
calc-wer-sentence: $(test_dir)/recognized.f.txt $(test_dir)/ref.f.txt
	@echo "FILE = $(snt)"
	$(scripts_dir)/compute-wer --mode=present --verbose=10 \
	ark:$(test_dir)/ref.f.txt "ark:cat $(test_dir)/recognized.f.txt | grep $(snt) |" 
$(test_dir)/wer-every.txt: $(test_dir)/wer.txt
	ls -1 $(test_dir)/result/*.txt | xargs -n 1 -I {} sh -c "f=\$$(basename {} | sed 's|.txt||'	);echo \$$f; $(MAKE) calc-wer-sentence snt=\$$f" \
		| grep -E "WER|FILE" | paste - - -d- > $@
		# | grep -E "WER" | awk '{print $$2}' > $@		

calc-wer-by-one: $(test_dir)/wer-every.txt
	cat $(test_dir)/wer-every.txt
#####################################################################################
calc-wer-kaldi-sentence: $(test_dir)/kaldi_rec.f.txt $(test_dir)/ref.f.txt
	@echo "FILE = $(snt)"
	$(scripts_dir)/compute-wer --mode=present --verbose=10 \
		ark:$(test_dir)/ref.f.txt  "ark:cat $(test_dir)/kaldi_rec.f.txt | grep $(snt) |"
calc-wer-kaldi-by-one: $(test_dir)/wer.txt
	ls -1 $(test_dir)/result/*.txt | xargs -n 1 -I {} sh -c "f=\$$(basename {} | sed 's|.txt||'	);echo \$$f; $(MAKE) calc-wer-kaldi-sentence snt=\$$f" \
		| grep -E "WER|FILE" | paste - - -d-
#####################################################################################
info: 
	@echo "UPLOAD		: $(UPLOAD_SERVICE)"
	@echo "STATUS		: $(STATUS_SERVICE)"
	@echo "RESULT		: $(RESULT_SERVICE)"
	@echo "Dir		: $(test_dir)"
	@echo "Model		: $(model)"
	@echo "Model med	: $(model_med)"
	
#####################################################################################
clean: clean-upload
	rm -rf $(test_dir)

clean-upload:
	rm -rf $(test_dir)/.upload.done

clean-wer-every:
	rm -rf  $(test_dir)/wer-every.txt 

clean-wer:
	rm -rf  $(test_dir)/wer.txt 	

clean-result:
	rm -rf $(test_dir)/result
	rm -rf $(test_dir)/recognized.txt 

clean-calc:
	rm -rf $(test_dir)/result/fixed
	rm -rf $(test_dir)/recognized.txt 

.PHONY:
	transcription-status info calc-wer
