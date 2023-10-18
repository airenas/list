import argparse
import os
import sys

from scripts.logger import logger
from scripts.rttm_to_seg import Seg

# join sil to speech if less than
join_sil = 0.5
max_speaker_speech = 30
gap_same_speaker_join = 3


def join_speaker(segs):
    segs = sorted(segs, key=lambda d: (d.start, d.sp))
    res = []
    for s in segs:
        if len(res):
            previous = res[-1]
            if previous.sp == s.sp and s.end - previous.start < max_speaker_speech \
                    and s.start - previous.end < gap_same_speaker_join:
                previous.end = s.end
                continue
        res.append(s)
    return res


def drop_silence(segs, file_len):
    segs = sorted(segs, key=lambda d: (d.start, d.sp))
    previous = None
    for s in segs:
        if previous:
            gap = s.start - previous.end
            take_sil = min(gap / 2, join_sil)
            if take_sil > 0:
                previous.end += take_sil
                s.start -= take_sil
        else:
            gap = s.start - 0
            take_sil = min(gap, join_sil)
            if take_sil > 0:
                s.start -= take_sil
        previous = s
    # add last
    if previous and previous.end < file_len:
        take_sil = min(file_len - previous.end, join_sil)
        if take_sil > 0:
            s.end += take_sil
    return segs


def main(argv):
    logger.info("Starting")
    parser = argparse.ArgumentParser(description="Drops silences, joins short speeches")
    parser.add_argument("--input", nargs='?', required=True, help="Input file to parse")
    parser.add_argument("--output", nargs='?', required=True, help="Output file")
    parser.add_argument("--len", nargs='?', type=float, required=True, help="Audio len")
    args = parser.parse_args(args=argv)

    file = args.input
    out_file = args.output
    file_name, _ = os.path.splitext(os.path.basename(out_file))
    logger.info(f"In: {file}")
    logger.info(f"Out: {out_file}")
    logger.info(f"Len: {args.len}")

    segs = []
    with open(file, "r") as f:
        for line in f:
            line = line.strip()
            if line:
                splits = line.split(" ")
                segs.append(Seg(start=splits[3], dur=splits[4], sp=splits[7]))

    segs = drop_silence(segs)
    segs = join_speaker(segs)

    for s in segs:
        s.dur = s.start - s.end

    rttm_lines = []
    for s in segs:
        rttm_line = f"SPEAKER file 1 {s.start:.3f} {s.dur:.3f} <NA> <NA> {s.sp} <NA> <NA>"
        rttm_lines.append(rttm_line)
    with open(out_file, "w") as file:
        [file.write(line + '\n') for line in rttm_lines]
    logger.info("done")


if __name__ == "__main__":
    main(sys.argv[1:])
