# The directory structure for Punctuation Model

```bash
<punct-model-dir>/
    settings.yml    # model settings,
                    #   punctuation vocaburary, etc...
    vocabulary      # words vocabulary (zero base index)
```

## Tensorflow model structure

```bash
<tf-model-dir>/
    <model-name>/           # model dir
        <model-version>/    # model version dir
            assets/         # tf dir
            variables/      # tf dir
            saved_model.pb  # ft files
```
