# Diegimas vienoje mašinoje naudojant Docker

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
| RAM  | 16 Gb | 32 Gb | |

Operacinė sistema: Linux OC 64-bit. Papildomai turi būti sudiegta:

| Komponentas | Min versija | URL |
| ---|-|-|
| Docker | 18.09.7 | [Link](https://docs.docker.com/engine/install/)
| Docker-compose | 1.23.0 | [Link](https://docs.docker.com/compose/install/) |
| GNU Make | 4.1 | [Link](https://www.gnu.org/software/make/manual/make.html) |
| Git | 2.17.1 | [Link](https://git-scm.com/download/linux) |

## Diegimas

1. Parsisiųskite diegimo skriptus (Ši git repositorija):

```bash
git clone https://bitbucket.org/airenas/list.git
cd list/deploy/run-docker
```

Docker diegimo direktorija yra *list/deploy/run-docker*.

2. Paruoškite konfigūracinį diegimo failą *Makefile.options*:

```bash
cp Makefile.options.template Makefile.options
```

3. Sukonfigūruokite *Makefile.options*:

| Parametras | Paskirtis | Pvz |
| ---|-|-|
| *deply_dir* | Pilnas kelias iki instaliavimo direktorijos | /home/user/list
| *models* | Instaliuojami modeliai. Galimi pasirinkimai: | |
| rabbitmq_pass | Eilės serviso slaptažodis ||
| mongo_pass | DB slaptažodis ||
| http_port | HTTP portas, kuriuo bus pasiekiami servisai | 80 |
| https_port | HTTPS portas, kuriuo bus pasiekiami servisai | 443 |
| host_external_url | Kompiuterio url, kuriuo servisai pasiekiami iš išorės. Naudojama nuorodai el. laiške | https://airenas.eu:7054 |
| smtp_host | SMTP serveris, laiškų siuntimui | 80 |
| smtp_port | SMTP portas | 587 |
| smtp_username | SMTP serverio vartotojas | olia@gmail.com |
| smtp_password | SMTP slaptažodis |  |

4. Instaliuokite

```bash
make install
```

Skriptas parsiųs reikalingus failus, paleis docker konteinerius. Priklausomai nuo inteneto ryšio diegimas gali užtrukti nuo 30 min iki kelių valandų.
Sistema bus sudiegta <deploy_dir> direktorijoje

## Patikrinimas

Atidarykite URL naršyklėje: *<host_external_url>/ausis/*. Turi atsidaryti puslapis.

Patikrinkite ar visi servisai veikia su docker-compose:

```bash
    cd <deploy_dir>
    docker-compose ps
```

Visi servisai turi būti *Up* būsenoje.

## Servisų sustabdymmas/valdymas

Servisai valdomi su docker-compose komanda:

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
