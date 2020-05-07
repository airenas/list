# The directory structure for Language Model

```bash
<lm-dir>/
    G.fst               # main lm fst
    words.txt           # vocabulary  
```

## RNNLM used for rescoring

```bash
<rnn-lm-dir>/
    feat_embedding.final.mat  # feature embedding matrix
    final.raw                 # rnn data
    info.txt                  # rnn load information data
    special_symbol_opts.txt   # special configuration for rnnlm
                              # related to original <lm-dir>/words.txt
    word_feats.txt            # word feature descriptors  
```
