# make file to join results from several files
-include Makefile.options

OUTPUT_ROOT?=/decoded
RESULTS_ROOT?=/results
APP_DIR?=/app

###############################################################################
LATTICES_TO_TEXT_APP?=$(APP_DIR)/lattices.to.text
LATTICES_TO_WEBVTT_APP?=$(APP_DIR)/lattices.to.webvtt
SEND_METRIC_APP?=${APP_DIR}/send.metric
###############################################################################

PATH:=utils:$(APPS_ROOT)/bin:$(PATH)

log=echo $$(date +'%Y-%m-%d %T')
logb=$(log): START
loge=$(log): END
###############################################################################
## for metrics
SHELL=$(APP_DIR)/m_shell.sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
worker=join.result
###############################################################################
# init symbolic directories
.SECONDARY:
.DELETE_ON_ERROR:
###############################################################################
$(RESULTS_ROOT)/$(id): 
	mkdir -p $@
###############################################################################
#$(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt
input_lat_files=$(patsubst %,$(OUTPUT_ROOT)/trans/%/punctuated/lat.restored.txt,$(CH_IDS))
###############################################################################
$(RESULTS_ROOT)/$(id)/resultFinal.txt: $(input_lat_files) | $(RESULTS_ROOT)/$(id)
	tsk=joinLattice;id=$(id);. $(APP_DIR)/m_start.sh $@
	cat $(input_lat_files) > $@

$(RESULTS_ROOT)/$(id)/webvtt.txt: $(input_lat_files) | $(RESULTS_ROOT)/$(id)
	tsk=lattice2webvtt;id=$(id);. $(APP_DIR)/m_start.sh $@
	$(LATTICES_TO_WEBVTT_APP) --namesMap "$(CH_IDS_FN)" $(input_lat_files) > $@
###############################################################################
info: 
	echo "ids = " $(CH_IDS)
	echo "ids_fn = " $(CH_IDS_FN)
	echo "input_lat_files = " $(input_lat_files)
###############################################################################
$(RESULTS_ROOT)/$(id)/.done: info $(RESULTS_ROOT)/$(id)/resultFinal.txt \
			$(RESULTS_ROOT)/$(id)/webvtt.txt
	touch $@

.PHONY: info	
