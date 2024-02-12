import argparse
import sys
import time

import torch

from scripts.logger import logger


def main(argv):
    logger.info("Starting")
    parser = argparse.ArgumentParser(description="VAD")
    parser.add_argument("--input", nargs='?', required=True, help="wav file")
    parser.add_argument("--output", nargs='?', required=True, help="out file")
    parser.add_argument("--rate", nargs='?', type=int, default=16000, help="audio rate")
    args = parser.parse_args(args=argv)

    torch.set_num_threads(1)
    logger.info("Init models")
    model, utils = torch.hub.load(repo_or_dir='snakers4/silero-vad', model='silero_vad', force_reload=False)
    (get_speech_timestamps, _, read_audio, *_) = utils
    sampling_rate = args.rate

    if sampling_rate != 8000 and sampling_rate != 16000:
        raise RuntimeError(f"wrong sampling rate {sampling_rate}. Wanted 8000 or 16000")

    wav = read_audio(args.input, sampling_rate=sampling_rate)

    if len(wav) == 0:
        raise RuntimeError(f"no wav in {args.input} ??")

    logger.info(f"Starting vad: file {args.input}'")
    start_time = time.time()
    speech_timestamps = get_speech_timestamps(wav, model, sampling_rate=sampling_rate)
    end_time = time.time()
    elapsed_time = end_time - start_time
    rt = elapsed_time / (len(wav) / sampling_rate)
    logger.info(f"Done vad in {elapsed_time:.2f}s, rt {rt}")

    rttm_lines = []
    for r in speech_timestamps:
        start = r["start"]
        dur = r["end"] - r["start"]
        start, dur = start / sampling_rate, dur / sampling_rate
        rttm_line = f"SPEAKER audio 1 {start:.3f} {dur:.3f} <NA> <NA> speaker <NA> <NA>"
        rttm_lines.append(rttm_line)

    with open(args.output, "w") as f:
        for r in rttm_lines:
            f.write(r + "\n")


if __name__ == "__main__":
    main(sys.argv[1:])
