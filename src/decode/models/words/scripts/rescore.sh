#!/bin/bash

# Copyright 2019 VDU (Author: Airenas Vaičiūnas).
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

modeldir=$1
inputfile=$2
outputfile=$3
## execute
echo "============= execute real rescore script ==================="
lattice-scale --inv-acoustic-scale=12.0 "ark:gunzip -c $inputfile|" ark:- | \
  lattice-add-penalty --word-ins-penalty=3.0 ark:- "ark:| gzip -c > $outputfile" || exit 1
echo "============= done ==================="
exit 0;
