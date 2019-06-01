## select all *.wav files from wav dir and pass to transcriber
## the result is *.wav.txt

ls -1 wav/*.wav | xargs -n1 -P6 ./tr.sh