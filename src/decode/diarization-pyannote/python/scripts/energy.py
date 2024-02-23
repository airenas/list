import argparse
import sys
import time

import librosa
import numpy as np

from scripts.logger import logger


def calculate_intensity(audio_file, frame_size_ms=25, shift_ms=10):
    y, sr = librosa.load(audio_file, sr=None)
    frame_size = int(sr * (frame_size_ms / 1000))
    shift = int(sr * (shift_ms / 1000))

    intensity = []
    for i in range(0, len(y) - frame_size + 1, shift):
        frame = y[i:i + frame_size]
        intensity.append(np.mean(frame ** 2))

    return intensity


def main(argv):
    logger.info("Starting")
    parser = argparse.ArgumentParser(description="Calculate energy")
    parser.add_argument("--input", nargs='?', required=True, help="wav file")
    parser.add_argument("--output", nargs='?', required=True, help="out file")
    # parser.add_argument("--rate", nargs='?', type=int, default=16000, help="audio rate")
    args = parser.parse_args(args=argv)
    logger.info(f"Starting energy calculation: file {args.input}'")
    start_time = time.time()
    intensity = calculate_intensity(audio_file=args.input)
    end_time = time.time()
    elapsed_time = end_time - start_time
    logger.info(f"Done energy calc in {elapsed_time:.2f}s")

    with open(args.output, "w") as f:
        for v in intensity:
            f.write(f"{v}\n")


if __name__ == "__main__":
    main(sys.argv[1:])
