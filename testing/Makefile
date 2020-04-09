-include Makefile.options
-include $(config)
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
$(test_dir)/result:
	mkdir -p $@

$(test_dir)/result/.done: $(test_dir)/file.list | $(test_dir)/result
	cat $(test_dir)/file.list | xargs -n2 -P20 scripts/result.sh $(test_dir)/result
	touch $(test_dir)/result/.done
get_result: $(test_dir)/result/.done	
#####################################################################################
$(test_dir)/recognized.txt: $(test_dir)/result/.done
	awk '{print FILENAME,$$0}' $(test_dir)/result/*.txt | \
		sed 's|$(test_dir)/result/||' > $(test_dir)/recognized.txt

$(test_dir)/wer.txt: $(test_dir)/recognized.txt
	scripts/compute-wer --mode=present \
		ark:$(test_dir)/testdata/txt/ref.txt  ark:$(test_dir)/recognized.txt > $(test_dir)/wer.txt
calc_wer: $(test_dir)/wer.txt
	cat $(test_dir)/wer.txt
compare: $(test_dir)/recognized.txt	
	vimdiff $(test_dir)/testdata/txt/ref.txt $(test_dir)/recognized.txt
#####################################################################################
info: 
	@echo "UPLOAD: $(UPLOAD_SERVICE)"
	@echo "STATUS: $(STATUS_SERVICE)"
	@echo "RESULT: $(RESULT_SERVICE)"
	@echo "Dir:    $(test_dir)"
	
#####################################################################################
clean: clean_send
	rm -rf $(test_dir)

clean_upload:
	rm -rf $(test_dir)/.upload.done

clean_result:
	rm -rf $(test_dir)/result
	rm -rf $(test_dir)/recognized.txt 

.PHONY:
	transcription_status info calc_wer


	
