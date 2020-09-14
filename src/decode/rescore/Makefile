# make file to rescore lattices
-include Makefile.options

OUTPUT_ROOT?=/decoded
APPS_ROOT?=/kaldi
FINAL_PASS=2-pass

###############################################################################
SEND_METRIC_APP?=${APP_DIR}/send.metric
###############################################################################

PATH:=utils:$(APPS_ROOT)/bin:$(PATH)

# init symbolic directories
steps:
	ln -s $(APPS_ROOT)/steps

utils:
	ln -s $(APPS_ROOT)/utils

.SECONDARY:
.DELETE_ON_ERROR:	

###############################################################################
## for metrics
SHELL=$(APP_DIR)/m_shell.sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
worker=rescore
###############################################################################

#rescore lattice
$(OUTPUT_ROOT)/trans/%/$(FINAL_PASS)/.done: $(INPUT_ROOT)/trans/%/1-pass/decode/lat.1.gz | steps utils
	tsk=rescore;id=$*;. $(APP_DIR)/m_start.sh $@

	mkdir -p $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode
	## run rescore script for model
	LM_DIR=$(LM_DIR) $(SCRIPTS_DIR)/rescore.sh $(MODELS_ROOT) $< $(OUTPUT_ROOT)/trans/$*/$(FINAL_PASS)/decode/lat.1.gz || exit 1
	touch $@ 
