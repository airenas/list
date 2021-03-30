import argparse
import sys
from collections import defaultdict
import os

def main(argv):
    parser = argparse.ArgumentParser(description="Counts words, adds </s> as a word for each line",
                                     epilog="E.g. cat input.txt | " + sys.argv[0] + " > result.txt",
                                     formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    parser.add_argument("--files", default='', type=str, help="File list", required=True)
    parser.add_argument("--ref", default='', type=str, help="Ref file", required=True)
    args = parser.parse_args(args=argv)

    print("Starting", file=sys.stderr)

    file_index = defaultdict(str)
    ind = 1
    with open(args.files, 'r') as file:
        for line in file:
            # de7b16c9-f5bd-4e9d-a03e-a6de5b4f00dd 20210329-ben-tel-intelektika-1/wav/IPREK1_03.wav
            line = line.strip()
            strs = line.split()
            f = os.path.basename(strs[1])
            f_name, f_ext = os.path.splitext(f)
            file_index[f_name] = str(ind)
            ind = ind + 1

    with open(args.ref, 'r') as file:
        for line in file:
            line = line.strip()
            strs = line.split()
            v = strs[0]
            strs[0] = file_index[v]
            if (strs[0] == ""):
                raise ValueError('No ' + v)
            print(" ".join(strs))

    print("Done", file=sys.stderr)

if __name__ == "__main__":
    main(sys.argv[1:])