#!/bin/bash

# Copyright 2020 VDU (Author: Airenas Vaičiūnas).
# BSD-3

# This script preloads lm rescore model
# Begin configuration section.
PATH=utils:/apps/kaldi/bin:${PATH}

# Init pipes
if [ ! -p pipe_input ]; then
    mkfifo pipe_input
fi
if [ ! -p pipe_output ]; then
    mkfifo pipe_output
fi

if [ -z "${LM_DIR}" ]; then
  echo "NO LM_DIR env variable!"
  exit 1
fi
if [ -z "${RNNLM_DIR}" ]; then
  echo "NO RNNLM_DIR env variable!"
  exit 1
fi

# Run preload cmd
oldlm=$LM_DIR/G.fst
special_symbol_opts=$(cat $RNNLM_DIR/special_symbol_opts.txt)
word_embedding="rnnlm-get-word-embedding $RNNLM_DIR/word_feats.txt $RNNLM_DIR/feat_embedding.final.mat -|"

lattice-lmrescore-kaldi-rnnlm-pruned-pipe --lm-scale=0.45 $special_symbol_opts \
    --lattice-compose-beam=4 --acoustic-scale=0.1 --max-ngram-order=4 --normalize-probs=false \
    --use-const-arpa=false $oldlm "$word_embedding" "$RNNLM_DIR/final.raw" pipe_input
    