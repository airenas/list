#!/bin/bash

# Copyright 2020 VDU (Author: Airenas Vaičiūnas).
# BSD-3

# This script preloads decoder model
# Begin configuration section.
PATH=utils:/apps/kaldi/bin:${PATH}

# Init pipes
if [ ! -p pipe_input ]; then
    mkfifo pipe_input
fi
if [ ! -p pipe_output ]; then
    mkfifo pipe_output
fi

if [ -z "${MODELS_ROOT}" ]; then
  echo "NO MODELS_ROOT env variable!"
  exit 1
fi

# Run preload cmd
nnet3-latgen-faster-parallel-pipe \
     --num-threads=4 \
     --online-ivector-period=10 \
     --frame-subsampling-factor=3 \
     --frames-per-chunk=90 \
     --extra-left-context=40 \
     --extra-right-context=40 \
     --extra-left-context-initial=0 \
     --extra-right-context-final=0 \
     --minimize=false --max-active=7000 --min-active=200 --beam=15 \
     --lattice-beam=8 --acoustic-scale=1.0 --allow-partial=true \
     --word-symbol-table=$MODELS_ROOT/words.txt $MODELS_ROOT/final.mdl \
     $MODELS_ROOT/HCLG.fst pipe_input
