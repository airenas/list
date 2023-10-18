import argparse
import os
import sys
import time

import torch
from pyannote.audio import Pipeline

from scripts.logger import logger


def main(argv):
    logger.info("Starting")
    parser = argparse.ArgumentParser(description="Pyannote diarization")
    parser.add_argument("--input", nargs='?', required=True, help="wav file")
    parser.add_argument("--output", nargs='?', required=True, help="rttm file")
    parser.add_argument("--num-speakers", nargs='?', type=int, required=False, help="speakers count")
    args = parser.parse_args(args=argv)

    logger.info("Init models")
    pipeline = Pipeline.from_pretrained("pyannote/speaker-diarization@2.1", use_auth_token=os.getenv('HF_API_TOKEN'))
    cuda = os.getenv('CUDA')
    if cuda and cuda != "cpu":
        pipeline = pipeline.to(torch.device(cuda))
    logger.info(f"Starting diarization on '{cuda}'")
    logger.info(f"Num speakers {args.num_speakers}")
    start_time = time.time()
    diarization = pipeline(args.input, num_speakers=args.num_speakers)
    end_time = time.time()
    elapsed_time = end_time - start_time
    logger.info(f"Done diarization in {elapsed_time:.2f}s")

    with open(args.output, "w") as f:
        diarization.write_rttm(f)


if __name__ == "__main__":
    main(sys.argv[1:])
