#!/bin/bash

set -e

awk '
    function trim(file) {
        sub(".*/", "", file)
        sub(".wav.txt", ".txt", file)
        return file
    }
    {print trim(FILENAME),$0}' testdata/wav/t*.txt > recognized.txt

../scripts/compute-wer --mode=present ark:testdata/txt/ref.txt  ark:recognized.txt

exit 0;
