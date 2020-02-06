#!/bin/bash

# Copyright 2020 VDU (Author: Airenas Vaičiūnas).
# BSD-3

# This script does rescoring for words model

# Begin configuration section.
# End configuration section.
echo "$0 $@"  # Print the command line for logging

if [ $# -ne 3 ]; then
  echo "Usage: $0 <models-dir> <input-latice-file> <output-latice-file>"
  echo "e.g.:   model_dir/scripts/rescore.sh model_dir lat.1.gz lat.2.gz"
  exit 1;
fi

model_dir=$1
inputfile=$2
outputfile=$3

if [ -z "${LM_DIR}" ]; then
  echo "NO LM_DIR env variable!"
  exit 1
fi
if [ -z "${RNNLM_DIR}" ]; then
  echo "NO RNNLM_DIR env variable!"
  exit 1
fi
## execute
echo "============= execute real rescore script ==================="
echo "============= rnnlm rescore ================================="
oldlm=$LM_DIR/G.fst
special_symbol_opts=$(cat $RNNLM_DIR/special_symbol_opts.txt)
word_embedding="rnnlm-get-word-embedding $RNNLM_DIR/word_feats.txt $RNNLM_DIR/feat_embedding.final.mat -|"

lattice-lmrescore-kaldi-rnnlm-pruned --lm-scale=0.45 $special_symbol_opts \
    --lattice-compose-beam=4 --acoustic-scale=0.1 --max-ngram-order=4 --normalize-probs=false \
    --use-const-arpa=false $oldlm "$word_embedding" "$RNNLM_DIR/final.raw" \
    "ark:gunzip -c $inputfile|" "ark,t:|gzip -c>$outputfile.rnnlm_rescored" || exit 1;  

echo "============= scale rescore ================================="
lattice-scale --inv-acoustic-scale=12.0 "ark:gunzip -c $outputfile.rnnlm_rescored|" ark:- | \
  lattice-add-penalty --word-ins-penalty=3.0 ark:- "ark:| gzip -c > $outputfile" || exit 1
echo "============= done ==================="
exit 0;
