import argparse
import os
import sys

from scripts.logger import logger
from scripts.rttm_to_seg import Seg


class Params:
    def __init__(self):
        self.join_sil = 0.5
        self.max_speaker_speech = 30
        self.gap_same_speaker_join = 3


def join_speaker(segs, params):
    segs = sorted(segs, key=lambda d: (d.start, d.sp))
    res = []
    for s in segs:
        if len(res):
            previous = res[-1]
            if previous.sp == s.sp and s.end - previous.start < params.max_speaker_speech \
                    and s.start - previous.end < params.gap_same_speaker_join:
                previous.end = s.end
                continue
        res.append(s)
    return res


def drop_silence(segs, file_len, params):
    segs = sorted(segs, key=lambda d: (d.start, d.sp))
    previous = None
    for s in segs:
        if previous:
            gap = s.start - previous.end
            take_sil = min(gap / 2, params.join_sil)
            if take_sil > 0:
                previous.end += take_sil
                s.start -= take_sil
        else:
            gap = s.start - 0
            take_sil = min(gap, params.join_sil)
            if take_sil > 0:
                s.start -= take_sil
        previous = s
    # add last
    if previous and previous.end < file_len:
        take_sil = min(file_len - previous.end, params.join_sil)
        if take_sil > 0:
            s.end += take_sil
    return segs


def main(argv):
    logger.info("Starting")
    parser = argparse.ArgumentParser(description="Drops silences, joins short speeches")
    parser.add_argument("--input", nargs='?', required=True, help="Input file to parse")
    parser.add_argument("--output", nargs='?', required=True, help="Output file")
    parser.add_argument("--len", nargs='?', type=float, required=True, help="Audio len")
    parser.add_argument("--join-sil", nargs='?', type=float, default=0.5, help="join sil to speech if less than")
    parser.add_argument("--join-gap", nargs='?', type=float, default=3, help="join sil gap if same speaker")
    args = parser.parse_args(args=argv)

    file = args.input
    out_file = args.output
    file_name, _ = os.path.splitext(os.path.basename(out_file))
    logger.info(f"In: {file}")
    logger.info(f"Out: {out_file}")
    logger.info(f"Len: {args.len}")
    params = Params()
    params.join_sil = args.join_sil
    params.gap_same_speaker_join = args.join_gap

    logger.info(f"join sil: {params.join_sil}, gap: {params.gap_same_speaker_join}")

    segs = []
    with open(file, "r") as f:
        for line in f:
            line = line.strip()
            if line:
                splits = line.split(" ")
                segs.append(Seg(start=splits[3], dur=splits[4], sp=splits[7]))

    segs = drop_silence(segs, args.len, params)
    segs = join_speaker(segs, params)

    for s in segs:
        s.dur = s.end - s.start

    rttm_lines = []
    for s in segs:
        rttm_line = f"SPEAKER file 1 {s.start:.3f} {s.dur:.3f} <NA> <NA> {s.sp} <NA> <NA>"
        rttm_lines.append(rttm_line)
    with open(out_file, "w") as file:
        [file.write(line + '\n') for line in rttm_lines]
    logger.info("done")


if __name__ == "__main__":
    main(sys.argv[1:])
