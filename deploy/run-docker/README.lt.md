# Diegimas vienoje mašinoje naudojant *Docker*

## Apie

Transkribatroriaus IT sprendimas yra realizuotas *Docker* komponentais. Visa sistema sukonfigūruota ir paruošta paleisti su *docker-compose* konfigūraciniu failu. Sistemos darbui taip pat reikalingi kai kurie papildomi (nedokerizuoti) binariniai vykdomieji failai ir lietuvių kalbos atpažinimo modelių failai. Diegient Jums reikės:

- atlikti pakeitimus konfiguraciniame faile,

- nauojantis paruoštų skriptų pagalba, parsiųsti reikalingus binarinius ir modelių failus,
- su *docker-compose* startuoti sistemą.
Transcribatroriaus IT sprendimas gali būti diegiamas, bet kurioje operacinėje sistemoje, kuri palaiko Docker technologiją, bet buvo testuotas ir šis aprašymas apima tik *Linux* sistemas. *Win* ir *Mac* operacinėse sistemose tikėtina bus reikalingas papildomas *docker-compose.yml* failo pritaikymas.

## Reikalavimai

Aparatūrai:

| Komponentas | Min reikalavimai | Rekomenduojama | Papildomai |
| ---|-|-|-|
| Platform | x86_64 | | |
| CPU | 64-bit, 2 branduoliai | 8 branduoliai | |
| HDD | 40 Gb | | Priklausomai nuo sudiegtų atpažinimo modelių. Vienam modeliui papildomai reikia apie 10 Gb |
| RAM | 24 Gb | 32 Gb | |

Operacinė sistema: Linux OS 64-bit (papildomai žiūrėkite [reikalavimus Docker instaliacijai](https://docs.docker.com/engine/install/). Turi būti sudiegta:

| Komponentas | Min versija | URL |
| ---|-|-|
| Docker | 18.09.7 | [Link](https://docs.docker.com/engine/install/)
| Docker-compose | 1.23.0 | [Link](https://docs.docker.com/compose/install/) |

Papildomi įrankiai naudojami instaliuojant: [make](https://www.gnu.org/software/make/manual/make.html), [git](https://git-scm.com/download/linux), [wget](https://www.gnu.org/software/wget/manual/wget.html), [tar](https://www.gnu.org/software/tar/manual/).

## Prieš diegiant

Patikrinkite ar visi reikalingi komponentai veikia mašinoje:

```bash
    ## Docker
    docker run hello-world
    ## Docker-compose
    docker-compose --version
    ## Kiti komponentai
    make --version
    tar --version
    wget --version
    git --version
```

## Diegimas

1. Parsisiųskite diegimo skriptus (ši git repositorija):

    `git clone https://bitbucket.org/airenas/list.git`

    `cd list/deploy/run-docker`

    Docker diegimo skriptai yra direktorijoje yra *list/deploy/run-docker*.

1. Paruoškite konfigūracinį diegimo failą *Makefile.options*:

    `cp Makefile.options.template Makefile.options`

1. Sukonfigūruokite *Makefile.options*:

    | Parametras | Privalomas | Paskirtis |Pvz |
    | ---|-|-|-|
    | *deploy_dir* | + | Pilnas kelias iki instaliavimo direktorijos mašinoje. Šioje direktorijoje bus atsiųsti modeliai, sukurtas pakatalogis darbiniams transkribatoriaus failams | /home/user/list
    | *models* | + | Instaliuojami modeliai. Galimi pasirinkimai: *ben*(bendrinis modelis),  *ben-tel* (bendrinis modelis, telefoninė kokybė), *adm* (administracinis), *adm-tel*, *tei* (teisinis), *tei-tel*, *med* (medicininis), *med-tel* | ben adm-tel |
    | rabbitmq_pass | + | Eilės serviso slaptažodis ||
    | mongo_pass | + |DB slaptažodis ||
    | http_port | + | HTTP portas, kuriuo bus pasiekiami servisai mašinoje | 80 |
    | host_external_url | - | Kompiuterio URL, kuriuo servisai pasiekiami iš išorės. Naudojama nuorodai el. laiške | <https://airenas.eu:7054> |
    | smtp_host | - |SMTP serveris, laiškų siuntimui | 80 |
    | smtp_port | - |SMTP portas | 587 |
    | smtp_username | - | SMTP serverio vartotojas | olia@gmail.com |
    | smtp_password | - | SMTP slaptažodis |  |

1. Instaliuokite

    `make install -j4`

    Skriptas parsiųs reikalingus failus, paleis *docker* konteinerius. Priklausomai nuo inteneto ryšio diegimas gali užtrukti nuo 30 min iki kelių valandų.
    Sistema bus sudiegta *<deploy_dir>* direktorijoje

## Patikrinimas

1. Patikrinkite ar visi servisai veikia su *docker-compose*: `make status`. Visi servisai turi būti *Up* būsenoje.

1. Patikrinkite ar servisas gali priimti užklausas: `make status-service`. Turi grąžinti užklausos kodą 200.

1. Atidarykite URL naršyklėje: *<host_external_url>/ausis/*. Turi atsidaryti puslapis.

## Servisų sustabdymmas/valdymas

Servisai valdomi su *docker-compose* komanda:

```bash
    cd <deploy_dir>
    ##Sustabdymas
    docker-compose stop
    ##Paleidimas
    docker-compose up -d
```

## Išinstaliavimas

```bash
    make clean
```

Komandą reikia vykdyti *admin* teisėmis, pvz.: `sudo sh -c 'make clean'`.
