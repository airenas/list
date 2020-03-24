-include Makefile.options
#####################################################################################
file?=1.wav
model?=general

upload: 
	curl -X POST -k $(UPLOAD_SERVICE)/upload -H 'Accept: application/json' -H 'content-type: multipart/form-data' \
	-F recognizer=$(model) -F numberOfSpeakers=2 \
	-F file=@$(file)

record:
	arecord -f cd > $(file)
	
