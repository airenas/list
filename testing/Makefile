-include Makefile.options
-include $(config)
wordfix_dir=../../lm-scripts/wordfix
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
$(test_dir)/.data.done: | $(test_dir)
	curl -k https://prn509.vdu.lt:7080/testdata/testdata-v2.zip -o $(test_dir)/testdata.zip
	cd $(test_dir) && unzip testdata.zip
	rm $(test_dir)/testdata.zip
	touch $(test_dir)/.data.done
prepare: $(test_dir)/.data.done
#####################################################################################
t1.1: $(test_dir)/.data.done
	curl -X POST -i -k $(UPLOAD_SERVICE)/upload -H 'Accept: application/json' \
		-H 'content-type: multipart/form-data' -F file=@$(test_dir)/testdata/wav/t_aft_11.wav
t1.2.1: $(test_dir)/.data.done
	curl -X POST -i -k $(UPLOAD_SERVICE)/upload -H 'Accept: application/json' \
		-H 'content-type: multipart/form-data'
t1.2.2: $(test_dir)/.data.done
	curl -X POST -i -k $(UPLOAD_SERVICE)/upload -H 'Accept: application/json' \
		-H 'content-type: multipart/form-data' -F recognizer=fake-model-for-tc-1.1.2 \
		-F file=@$(test_dir)/testdata/wav/t_aft_11.wav
t2.1.1: $(test_dir)/.data.done
	curl -X POST -i -k $(UPLOAD_SERVICE)/upload -H 'Accept: application/json' \
		-H 'content-type: multipart/form-data' -F file=@$(test_dir)/testdata/wav/t_aft_11.wav
_tr_status: $(test_dir)/.data.done
	echo "tr id = $(tr_id)"
	curl -X GET -i -k $(STATUS_SERVICE)/status/$(tr_id) -H "accept: application/json"
t2.1.2: _tr_status
t2.2.1: $(test_dir)/.data.done
	curl -X POST -i -k $(UPLOAD_SERVICE)/upload -H 'Accept: application/json' \
		-H 'content-type: multipart/form-data' -F file=@$(test_dir)/testdata/wav/t_aft_22.wav	
t2.2.2: _tr_status	
t3.1: $(test_dir)/.data.done | $(test_dir)/result
	echo "tr id = $(tr_id)"
	curl -X GET -k $(RESULT_SERVICE)/result/$(tr_id)/result.txt -o $(test_dir)/result/t_aft_11.txt
	cat $(test_dir)/result/t_aft_11.txt
	diff $(test_dir)/testdata/txt/t_aft_11.txt $(test_dir)/result/t_aft_11.txt
t3.2: $(test_dir)/.data.done | $(test_dir)/result
	echo "tr id = $(tr_id)"
	curl -X GET -k $(RESULT_SERVICE)/audio/$(tr_id) -o $(test_dir)/result/t_aft_11.wav
	play $(test_dir)/result/t_aft_11.wav
	diff $(test_dir)/testdata/wav/t_aft_11.wav $(test_dir)/result/t_aft_11.wav

t4.1.1: $(test_dir)/.data.done
	curl -X POST -i -k $(UPLOAD_SERVICE)/upload -H 'Accept: application/json' \
		-H 'content-type: multipart/form-data' -F file=@$(test_dir)/testdata/wav/t_aft_12.wav
t4.1.2: _tr_status	
t4.1.3: $(test_dir)/.data.done | $(test_dir)/result
	echo "tr id = $(tr_id)"
	curl -X GET -k $(RESULT_SERVICE)/result/$(tr_id)/resultFinal.txt -o $(test_dir)/result/t_aft_12.txt
	cat $(test_dir)/result/t_aft_12.txt
	diff $(test_dir)/testdata/txt/t_aft_12.txt $(test_dir)/result/t_aft_12.txt
t4.2.1: $(test_dir)/.data.done
	curl -X POST -i -k $(UPLOAD_SERVICE)/upload -H 'Accept: application/json' \
		-H 'content-type: multipart/form-data' -F file=@$(test_dir)/testdata/wav/t_aft_13.wav
t4.2.2: _tr_status	
t4.2.3: $(test_dir)/.data.done | $(test_dir)/result
	echo "tr id = $(tr_id)"
	curl -X GET -k $(RESULT_SERVICE)/result/$(tr_id)/result.txt -o $(test_dir)/result/t_aft_13.txt
	cat $(test_dir)/result/t_aft_13.txt
	diff $(test_dir)/testdata/txt/t_aft_13.txt $(test_dir)/result/t_aft_13.txt
#####################################################################################
$(test_dir)/.upload.done: $(test_dir)/.data.done
	@echo "Sending data to server $(UPLOAD_SERVICE)"
	export UPLOAD_SERVICE=$(UPLOAD_SERVICE); \
	export SKIP_NUM_JOIN=1; \
		ls -1 $(test_dir)/testdata/wav/t_aft_4*.wav | xargs -n1 -P10 scripts/send.sh $(model) > $(test_dir)/file.list; \
		ls -1 $(test_dir)/testdata/wav/t_aft_5*.wav | xargs -n1 -P10 scripts/send.sh $(model_med) >> $(test_dir)/file.list;
	touch $(test_dir)/.upload.done
	
upload-to-server: $(test_dir)/.upload.done
recognizers:
	curl -k $(UPLOAD_SERVICE)/recognizers | jq

#####################################################################################
transcription-status: $(test_dir)/file.list
	export STATUS_SERVICE=$(STATUS_SERVICE); cat $(test_dir)/file.list | xargs -n2 -P20 scripts/status.sh | sort	
#####################################################################################
$(test_dir)/result:
	mkdir -p $@

$(test_dir)/result/.done: $(test_dir)/file.list | $(test_dir)/result
	export RESULT_SERVICE=$(RESULT_SERVICE); cat $(test_dir)/file.list | xargs -n2 -P20 scripts/result.sh $(test_dir)/result
	touch $(test_dir)/result/.done
get_result: $(test_dir)/result/.done	
#####################################################################################
$(test_dir)/result/fixed:
	mkdir -p $@

$(test_dir)/recognized.txt: $(test_dir)/result/.done | $(test_dir)/result/fixed
	ls -1 $(test_dir)/result/*$(tf).txt | xargs -n 1 -I {} sh -c "f={};\
		cat \$$f | tr '\n' ' ' | sed -r 's/ +/ /g' | sed -r 's/^ //' > $(test_dir)/result/fixed/\$${f##*/}"
	awk '{print FILENAME,$$0}' $(test_dir)/result/fixed/*.txt | \
		sed 's|$(test_dir)/result/fixed/||' > $(test_dir)/recognized.txt

$(test_dir)/recognized.f.txt: $(test_dir)/recognized.txt
		cat $(test_dir)/recognized.txt | perl $(wordfix_dir)/filter.pl $(wordfix_dir) > $(test_dir)/recognized.f.txt

$(test_dir)/wer.txt: $(test_dir)/recognized.f.txt
	scripts/compute-wer --mode=present \
		ark:$(test_dir)/testdata/txt/ref.txt  ark:$(test_dir)/recognized.f.txt > $(test_dir)/wer.txt
$(test_dir)/ref.txt: $(test_dir)/.data.done
	cat $(test_dir)/testdata/txt/ref.txt | sed 's/ *$$//g' > $@
calc-wer: $(test_dir)/wer.txt
	cat $(test_dir)/wer.txt
compare: $(test_dir)/recognized.f.txt	$(test_dir)/ref.txt
	vimdiff $(test_dir)/ref.txt $(test_dir)/recognized.f.txt
calc-wer-sentence: $(test_dir)/recognized.f.txt
	scripts/compute-wer --mode=present --verbose=10 \
		ark:$(test_dir)/testdata/txt/ref.txt  "ark:cat $(test_dir)/recognized.f.txt | grep $(snt) |"
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

clean-result:
	rm -rf $(test_dir)/result
	rm -rf $(test_dir)/recognized.txt 

clean-calc:
	rm -rf $(test_dir)/result/fixed
	rm -rf $(test_dir)/recognized.txt 

.PHONY:
	transcription-status info calc-wer


	
