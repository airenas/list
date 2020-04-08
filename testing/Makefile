-include Makefile.options
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
	curl -k https://prn509.vdu.lt:7080/testdata/testdata.zip -o $(test_dir)/testdata.zip
	cd $(test_dir) && unzip testdata.zip
	rm $(test_dir)/testdata.zip
	touch $(test_dir)/.data.done
prepare: $(test_dir)/.data.done
#####################################################################################
$(test_dir)/.upload.done: $(test_dir)/.data.done
	@echo "Sending data to server $(UPLOAD_SERVICE)"
	ls -1 $(test_dir)/testdata/wav/t_aft_4*.wav | xargs -n1 -P10 scripts/send.sh $(model) > $(test_dir)/file.list
	touch $(test_dir)/.upload.done
	
upload_to_server: $(test_dir)/.upload.done
#####################################################################################
transcription_status: $(test_dir)/file.list
	cat $(test_dir)/file.list | xargs -n2 -P20 scripts/status.sh | sort	
#####################################################################################
get_results: $(test_dir)/file.list
	cat $(test_dir)/file.list | xargs -n2 -P20 scripts/result.sh
#####################################################################################
calc_wer: $(test_dir)/file.list
	scripts/test_wer.sh 
#####################################################################################
info: 
	@echo "UPLOAD: $(UPLOAD_SERVICE)"
	@echo "STATUS: $(STATUS_SERVICE)"
	@echo "RESULT: $(RESULT_SERVICE)"
	
#####################################################################################

clean: clean_send
	rm -rf $(test_dir)

clean_upload:
	rm -rf $(test_dir)/.upload.done

.PHONY:
	transcription_status info


	
