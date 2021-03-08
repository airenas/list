#!/bin/bash

# Copyright 2021 VDU (Author: Airenas Vaičiūnas).
# BSD-3

# Sample restore script
set -e
echo "$0 $@"  # Print the command line for logging

if [ $# -ne 3 ]; then
  echo "Usage: $0 <models-dir> <input-latice-file> <output-latice-dir>"
  echo "e.g.:   model_dir/scripts/restore.sh model_dir lat.1.gz restore"
  exit 1;
fi

modeldir=$1
inputfile=$2
outputdir=$3
## execute
echo "============= execute real restore script ==================="

unk_id=$(grep '<unk>' ${modeldir}/words.txt | cut -d' ' -f2);
sil_id=$(grep '<eps>' ${modeldir}/words.txt | cut -d' ' -f2);
lattice-prune --beam=7 "ark:gunzip -c ${inputfile}|"  ark:- | \
	lattice-push ark:- ark:- | \
	lattice-align-words --silence-label=${sil_id} --partial-word-label=${unk_id} ${modeldir}/word_boundary.int \
	${modeldir}/final.mdl ark:- ark,t:${outputdir}/L1.lat || exit 1

# Extract 1-best lattice
# negalima naudoti lattice-to-nbest --n=1 ark,t:L1.lat ark,t:L2.lat
# lattice-to-nbest ne tik papildo '-1' utt_id, bet ir prijungia tylas prie þodþiø pabaigos
# taip desinchronizuodama L1 ir L2
lattice-1best ark,t:${outputdir}/L1.lat \
	ark,t:${outputdir}/L2.lat || exit 1

# Replace transition-ids by phone-ids (timing is lost)
lattice-to-phone-lattice --replace-words=false ${modeldir}/final.mdl \
	ark,t:${outputdir}/L1.lat \
	ark,t:${outputdir}/L3.lat  || exit 1

numJoinParam='--join-num 0.03'; if [ "${SKIP_NUM_JOIN}" == "1" ]; then numJoinParam=''; fi
spkJoinParam='--join-spk'; if [ "${SKIP_SPEAKER_JOIN}" == "true" ] || [ "${SKIP_SPEAKER_JOIN}" == "1" ]; then spkJoinParam=''; fi
echo "SKIP_NUM_JOIN=${SKIP_NUM_JOIN}; numJoinParam=${numJoinParam}"
echo "SKIP_SPEAKER_JOIN=${SKIP_SPEAKER_JOIN}; spkJoinParam=${spkJoinParam}"
# Perform processing 
cd restore && perl lat_restore.pl ${outputdir}/L1.lat \
	${outputdir}/L2.lat ${outputdir}/L3.lat ${modeldir}/words.txt \
	${modeldir}/phones.txt ${spkJoinParam} ${numJoinParam} > ${outputdir}/lat.restored.txt || exit 1

echo "============= done ==================="
exit 0;
