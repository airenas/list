import argparse
import os
import sys

from scripts.logger import logger


class Namer:
    def __init__(self):
        self.num = 0
        self.dict = {}

    def label(self, sp):
        res = self.dict.get(sp, None)
        if res:
            return res
        res = f"S{self.num:03d}"
        self.dict[sp] = res
        self.num += 1
        return res


class Seg:
    def __init__(self, start, dur, sp):
        self.sp = sp
        self.dur = float(dur)
        self.start = float(start)
        self.end = self.start + self.dur


def main(argv):
    logger.info("Starting")
    parser = argparse.ArgumentParser(description="Converts rttm to kaldi seg format")
    parser.add_argument("--input", nargs='?', required=True, help="Input file to parse")
    parser.add_argument("--output", nargs='?', required=True, help="Output file")
    args = parser.parse_args(args=argv)

    file = args.input
    out_file = args.output
    file_name, _ = os.path.splitext(os.path.basename(out_file))
    logger.info(f"In: {file}")
    logger.info(f"Out: {out_file}")

    segs = []
    with open(file, "r") as f:
        for line in f:
            line = line.strip()
            if line:
                splits = line.split(" ")
                segs.append(Seg(start=splits[3], dur=splits[4], sp=splits[7]))

    segs = sorted(segs, key=lambda d: (d.sp, d.start))

    rttm_lines = []
    namer = Namer()
    for s in segs:
        start_time = int(s.start * 100)
        duration = int(s.dur * 100)
        label = namer.label(s.sp)

        rttm_line = f"prepared 1 {start_time} {duration} NA NA NA {label}"
        rttm_lines.append(rttm_line)

    with open(out_file, "w") as file:
        [file.write(line + '\n') for line in rttm_lines]
    logger.info("done")


if __name__ == "__main__":
    main(sys.argv[1:])
