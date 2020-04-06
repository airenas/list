#!/bin/bash

# Copyright 2020 VDU (Author: Airenas Vaičiūnas).
# BSD-3

# This script does decoding for word model

# Begin configuration section.
nj=1 # number of decoding jobs.
num_threads=4 # if >1, will use gmm-latgen-faster-parallel
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
online_ivectors="'scp:$ivectorsdir/ivector_online.scp'"
feats="'ark,s,cs:apply-cmvn --norm-means=false --norm-vars=false --utt2spk=ark:$datadir/split$nj/1/utt2spk scp:$datadir/split$nj/1/cmvn.scp scp:$datadir/split$nj/1/feats.scp ark:- |'"
lat_wspecifier="'ark:|lattice-scale --acoustic-scale=10.0 ark:- ark:- | gzip -c >$decodedir/lat.1.gz'"

## execute decoding
echo "============= execute real decoding script ==================="
./pipe.runner -i pipe_input -o pipe_output -t 3m $online_ivectors $feats $lat_wspecifier || exit 1;
  # nnet3-latgen-faster-parallel \
  #    --num-threads=$num_threads \
  #    --online-ivectors=scp:"$ivectorsdir"/ivector_online.scp \
  #    --online-ivector-period=10 \
  #    --frame-subsampling-factor=3 \
  #    --frames-per-chunk=90 \
  #    --extra-left-context=40 \
  #    --extra-right-context=40 \
  #    --extra-left-context-initial=0 \
  #    --extra-right-context-final=0 \
  #    --minimize=false --max-active=7000 --min-active=200 --beam=15 \
  #    --lattice-beam=8 --acoustic-scale=1.0 --allow-partial=true \
  #    --word-symbol-table=$graphdir/words.txt $graphdir/final.mdl \
  #    $graphdir/HCLG.fst "$feats" "$lat_wspecifier" || exit 1;
echo "============= done ==================="
exit 0;
