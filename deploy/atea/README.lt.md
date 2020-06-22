# Diegimo instrukcija ATEA Kubernetes

Diegimo instrukcija ir pagalbiniai skriptai paruo�ti Linux OS.

## Prie� diegiant

- diegimo skriptai naudoja pagalbinius ?rankius *git*, *make*, *pwgen*, *rsync*, *kubectl*. ?sitikinkite, kad jie suinstaliuoti sistemoje.
- *kubectl* gali prisijungti prie ATEA k8s. Patikrinkite, kad  `kubectl get pods` komada veikia ir prisijungia prie ATEA k8s.

## Diegimas

1. Parsisi?skite diegimo skriptus (�i git repositorija):

    `git clone https://bitbucket.org/airenas/list.git`

    `cd list/deploy/atea`

    Docker diegimo skriptai yra direktorijoje yra *list/deploy/atea*.

1. Paruo�kite slapta�od�ius *rabbitmq* eil?s servisui ir mongo DB (*secrets/secrets.yml* failas):

    `make prepare-secrets`

1. Instaliuokite slapta�od�ius ? k8s:

    `make install-secrets`

1. I�trinkite *secrets.yml* lokal? fil�:

    `make clean-secrets`

1. Paruo�kite k8s saugyklas:

    `make install-volumes`

1. Sukopijuokite modelius ir pagalbinius binarinius failus ? k8s saugyklas:

    `make copy-data`

1. Instaliuokite servisus:

    `make install-services`

## Komponento veikimo patikrinimas

### *Upload* servisas

Paleid�iame proxy ? servis�: `make proxy-upload`

Kitame terminaliniame lange vykdome komand�: `make info-upload`

Turime gauti atsakym� su HTTP statuso kodu: 200.

### *Status* servisas

Analogi�kai, kaip *Upload* serviso patikrinimas: `make proxy-status`. Kitame terminaliniame lange vykdome komand�: `make info-status`

### *Result* servisas

Analogi�kai, kaip *Upload* serviso patikrinimas: `make proxy-result`. Kitame terminaliniame lange vykdome komand�: `make info-result`

## Prototipo servis? atnaujinimas

1. Padarome pakeitimus faile deployments/transcription.yml.

1. I�saugome pakeitimus git repozitorijoje.

1. Vykdome komand�: `make install-services`

## Duomen? atnaujinimas

1. Atnaujinus duomenis, bus pakeista ir �i repositorija su nuorodomis ? naujus modeli? failus. Patikrinkite, kad turite naujausius skriptus:

    `git pull`

1. Sukopijuokite atnaujintus modelius ir pagalbinius binarinius failus ? k8s saugyklas:

    `make copy-data`

## I�trynimas

1. I�trinkite servisus:

    `make clean-services`

1. I�trinkite duomenis:

    `make clean-data`

## Papildoma informacija

### Skriptai problem? sprendimui

**Servisas nestartuoja, informacija:** `kubectl describe deployment <deployment-name>`

**Serviso log:** `kubectl logs pod <pod-name>`

**Prisijungimas prie veikian�ios ma�inos:** `kubectl exec -it <pod-name> /bin/bash` arba `kubectl exec -it <pod-name> /bin/sh`.
