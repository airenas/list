import argparse
import os
import sys

from scripts.fix_short_silences import read_rttm, write_rttm, join_speaker
from scripts.logger import logger
from scripts.rttm_to_seg import Seg


class Params:
    def __init__(self):
        self.join_sil = 0.5
        self.max_speaker_speech = 30
        self.gap_same_speaker_join = 3
        self.len = 0


def find_sil(segments, i_seg, start):
    while i_seg < len(segments):
        seg = segments[i_seg]
        if seg.start > start:
            return i_seg
        i_seg += 1
    return i_seg


def end_pos(segments, i):
    if i < 0:
        return 0
    if i < len(segments):
        return segments[i].end
    if len(segments) > 0:
        return segments[-1].end
    return 0


def len_at_pos(segments, i):
    if i < 0 or i == len(segments):
        return 0, None
    return segments[i].end - segments[i].start, segments[i]


def start_pos(segments, i, a_len):
    if i < 0:
        return 0
    if i < len(segments):
        return segments[i].start
    return a_len


def is_space(segments, i, params):
    if i == len(segments):
        return True
    st = start_pos(segments, i, params.len)
    en = end_pos(segments, i - 1)
    return st > en


va_name = "ADDED-VA"


def add_va(segments, va_seqs, params):
    i = 0
    for s_va in va_seqs:
        i = find_sil(segments, i, s_va.start)
        ok = True
        while s_va.end > s_va.start and ok:
            ok = False
            if is_space(segments, i, params):
                st = max(s_va.start - params.join_sil, end_pos(segments, i - 1))
                en = min(s_va.end + params.join_sil, start_pos(segments, i, params.len))
                if i == len(segments) and en < s_va.end:
                    en = s_va.end
                if st < en:
                    n_seg = Seg(start=st, dur=en - st, sp=va_name)
                    segments.insert(i, n_seg)
                    ok = True
            s_va.start = min(s_va.end, start_pos(segments, i, params.len))
            i = find_sil(segments, i, s_va.start)
    return segments


def first_speaker(segments):
    for s in segments:
        if s.sp != va_name:
            return s.sp
    return "speaker_va"


def assign_speaker(segments, params):
    for i, s in enumerate(segments):
        if s.sp == va_name:
            p_len, p = len_at_pos(segments, i - 1)
            p_diff = s.start - end_pos(segments, i - 1)

            n_len, n = len_at_pos(segments, i + 1)
            n_diff = start_pos(segments, i + 1, params.len) - s.end
            if p_len > 0 and n_len > 0:
                if p_diff > n_diff:
                    s.sp = n.sp
                    continue
                if p_diff < n_diff and n.sp != va_name:
                    s.sp = p.sp
                    continue
            if p_len > n_len:
                s.sp = p.sp
                continue
            if n_len > 0 and n.sp != va_name:
                s.sp = n.sp
                continue
            if p:
                s.sp = p.sp
            else:
                s.sp = first_speaker(segments)
    return segments


def main(argv):
    logger.info("Starting")
    parser = argparse.ArgumentParser(description="Drops silences if there is a speed in vad data")
    parser.add_argument("--input", nargs='?', required=True, help="Input file to parse")
    parser.add_argument("--vad", nargs='?', required=True, help="VAD file in rttm format")
    parser.add_argument("--output", nargs='?', required=True, help="Output file")
    parser.add_argument("--len", nargs='?', type=float, required=True, help="Audio len")
    parser.add_argument("--join-sil", nargs='?', type=float, default=0.5, help="join sil to speech if less than")
    parser.add_argument("--join-gap", nargs='?', type=float, default=3, help="join sil gap if same speaker")
    args = parser.parse_args(args=argv)

    file = args.input
    vad_file = args.vad
    out_file = args.output
    file_name, _ = os.path.splitext(os.path.basename(out_file))
    logger.info(f"In: {file}")
    logger.info(f"VA: {vad_file}")
    logger.info(f"Out: {out_file}")
    logger.info(f"Len: {args.len}")
    params = Params()
    params.join_sil = args.join_sil
    params.gap_same_speaker_join = args.join_gap
    params.len = args.len

    logger.info(f"join sil: {params.join_sil}, gap: {params.gap_same_speaker_join}")

    segs = read_rttm(file)
    segs = sorted(segs, key=lambda d: (d.start, -d.end, d.sp))

    va_seqs = read_rttm(vad_file)
    va_seqs = sorted(va_seqs, key=lambda d: (d.start, -d.end, d.sp))

    segs = add_va(segs, va_seqs, params)
    segs = assign_speaker(segs, params)
    segs = join_speaker(segs, params)

    for s in segs:
        s.dur = s.end - s.start

    write_rttm(out_file, segs)
    logger.info("done")


if __name__ == "__main__":
    main(sys.argv[1:])
