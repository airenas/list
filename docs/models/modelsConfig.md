# The directory structure for Transcription models configuration

The configuration dir contains all transcription model configuration. Each model is configured in separate (.yml) file. File *recognizers.map.yml* contains initial recognizer mapping to a real model. It is used in [UploadService](https://app.swaggerhub.com/apis/aireno/Transkipcija/1.3.0#/upload/upload).

```bash
<config-dir>/
    regognizers.map.yml # recognizer name to model mapping
    <model name 1>.yml  # model 1 descriptor  
    ...
    <model name N>.yml  # model N descriptor
```

## File *recognizers.map.yml* sample

```yml
default: words_v3           # default model when recognizer=''. 
                            # real model config in <config-dir>/word_v3.yml
words_v3: words_v3          # some mapping words_v3 -> words_v3
general_standard: words_v3  # mapping general_standard -> words_v3
general_telephonic: words_t # mapping general_telephonic -> words_t
#           ...
```

## Sample model config *<model name 1>.yml*

```yml
name: AŽ v3 (RNNLM)
description: Žodžių atpažintuvas. Gailiaus paruoštas 2020 04 17. Su RNNLM
date_created: 2020-04-17
settings:
    # am model root dir
  models_root: /models/ac/20_04_17-2g/

  # scripts dir - where models decode.sh and rescore.sh are located
  scripts_dir: /models/config/words_v3/scripts/
  lm_dir: /models/lm/v3/
  rnnlm_dir: /models/rnnlm/v1/
  punctuate: true  # do need to punctuate result

  # transcription preload commands - may be shared between models if key is the same
  transcription-preload_key: words_v3
  transcription-preload_cmd: /models/config/words_v3/scripts/decode.preload.sh

  # transcription preload commands - may be shared between models if key is the same
  rescore-preload_key: words_v3
  rescore-preload_cmd: /models/config/words_v3/scripts/rescore.preload.sh
```
