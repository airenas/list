#!/bin/bash

# Copyright 2019 VDU (Author: Airenas Vaičiūnas).
# BSD-3

# This script does decoding for phone model

# Begin configuration section.
nj=4 # number of decoding jobs.
num_threads=1 # if >1, will use gmm-latgen-faster-parallel
cmd=run.pl
# End configuration section.

echo "$0 $@"  # Print the command line for logging

if [ $# -ne 4 ]; then
  echo "Usage: $0 <graph-dir> <ivector-dir> <data-dir> <decode-dir>"
  echo "e.g.:   model_dir/scripts/decode.sh model_dir ivector-dir data-dir output-dir"
  exit 1;
fi

graphdir=$1
ivectorsdir=$2
datadir=$3
decodedir=$4

## execute decoding
echo "============= execute real decoding script ==================="
steps/nnet3/decode.sh --num-threads $num_threads --nj $nj \
    --acwt 1.0  --post-decode-acwt 10.0 \
	--config conf/decode.conf \
	--skip-scoring true --cmd "$cmd" --nj 1 \
	--skip-diagnostics true \
	--extra-left-context 40  \
    --extra-right-context 40  \
    --extra-left-context-initial 0 \
    --extra-right-context-final 0 \
    --frames-per-chunk 90 \
	--online-ivector-dir $ivectorsdir \
    $graphdir $datadir $decodedir || exit 1;
echo "============= done ==================="
exit 0;
