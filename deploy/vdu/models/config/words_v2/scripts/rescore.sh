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

## execute
echo "============= execute real rescore script ==================="
echo "============= rnnlm rescore ================================="
outdir=$(dirname $outputfile)
mkdir -p $outdir/log
lat_rspecifier="'ark:gunzip -c $inputfile|'"
lat_wspecifier="'ark,t:|gzip -c>$outdir/rnnlm_rescored.gz'"
/app/pipe.runner -i pipe_input -o pipe_output -t 10m $lat_rspecifier $lat_wspecifier || exit 1;

echo "============= scale rescore ================================="
lattice-scale --inv-acoustic-scale=12.0 "ark:gunzip -c $outdir/rnnlm_rescored.gz|" ark:- | \
  lattice-add-penalty --word-ins-penalty=3.0 ark:- "ark:| gzip -c > $outputfile" || exit 1
echo "============= done ==================="
exit 0;
