# make file to rescore lattices
-include Makefile.options

OUTPUT_ROOT?=/decoded
APPS_ROOT?=/kaldi
FINAL_PASS=2-pass

PATH:=utils:$(APPS_ROOT)/bin:$(PATH)

log=@echo $$(date +'%Y-%m-%d %T')
logb=$(log): START
loge=$(log): END

# init symbolic directories
steps:
	ln -s $(APPS_ROOT)/steps

utils:
	ln -s $(APPS_ROOT)/utils

.SECONDARY:
.DELETE_ON_ERROR:	

#rescore lattice
$(OUTPUT_ROOT)/trans/%/$(FINAL_PASS)/.done: $(INPUT_ROOT)/trans/%/1-pass/decode/lat.1.gz | steps utils
	$(logb) $@
	mkdir -p $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode
	## run rescore script for model
	LM_DIR=$(LM_DIR) $(SCRIPTS_DIR)/rescore.sh $(MODELS_ROOT) $< $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode/lat.1.gz || exit 1
	touch $@ 
	$(loge) $@
