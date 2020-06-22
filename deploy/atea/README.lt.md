# Diegimo instrukcija ATEA Kubernetes

Diegimo instrukcija ir pagalbiniai skriptai paruošti Linux OS.

## Prieš diegiant

- diegimo skriptai naudoja pagalbinius įrankius *git*, *make*, *pwgen*, *rsync*, *kubectl*. Įsitikinkite, kad jie suinstaliuoti sistemoje.
- *kubectl* gali prisijungti prie ATEA k8s. Patikrinkite, kad  `kubectl get pods` komada veikia ir prisijungia prie ATEA k8s.

## Diegimas

1. Parsisiųskite diegimo skriptus (ši git repositorija):

    `git clone https://bitbucket.org/airenas/list.git`

    `cd list/deploy/atea`

    Docker diegimo skriptai yra direktorijoje yra *list/deploy/atea*.

1. Paruoškite slaptažodžius *rabbitmq* eilė servisui ir mongo DB (*secrets/secrets.yml* failas):

    `make prepare-secrets`

1. Instaliuokite slaptažodžius į k8s:

    `make install-secrets`

1. Ištrinkite *secrets.yml* lokalų filą:

    `make clean-secrets`

1. Paruoškite k8s saugyklas:

    `make install-volumes`

1. Sukopijuokite modelius ir pagalbinius binarinius failus į k8s saugyklas:

    `make copy-data`

1. Instaliuokite servisus:

    `make install-services`

## Komponento veikimo patikrinimas

### *Upload* servisas

Paleiskite proxy į servisą: `make proxy-upload`

Kitame terminaliniame lange vykdykite komandą: `make info-upload`

Turite gauti atsakymą su HTTP statuso kodu: 200.

### *Status* servisas

Analogiškai, kaip *Upload* serviso patikrinimas: `make proxy-status`. Kitame terminaliniame lange vykdome komandą: `make info-status`

### *Result* servisas

Analogiškai, kaip *Upload* serviso patikrinimas: `make proxy-result`. Kitame terminaliniame lange vykdome komandą: `make info-result`

## Servisų atnaujinimas

1. Padarome pakeitimus faile deployments/transcription.yml.

1. Išsaugome pakeitimus git repozitorijoje.

1. Vykdome komandą: `make install-services`

## Duomenų atnaujinimas

1. Atnaujinus duomenis, bus pakeista ir ši repositorija su nuorodomis į naujus modelių failus. Patikrinkite, kad turite naujausius skriptus:

    `git pull`

1. Sukopijuokite atnaujintus modelius ir pagalbinius binarinius failus į k8s saugyklas:

    `make copy-data`

## Ištrynimas

1. Ištrinkite servisus:

    `make clean-services`

1. Ištrinkite duomenis:

    `make clean-data`

1. Ištrinkite pagalbinį duomenų servisą:

    `make clean-vh`

## Papildoma informacija

### Skriptai problemų sprendimui

**Servisas nestartuoja, informacija:** `kubectl describe deployment <deployment-name>`

**Serviso log:** `kubectl logs pod <pod-name>`

**Prisijungimas prie veikiančios mašinos:** `kubectl exec -it <pod-name> /bin/bash` arba `kubectl exec -it <pod-name> /bin/sh`.