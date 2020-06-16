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

4. Instaliuokite

```bash
make install
```

Skriptas parsiųs reikalingus failus, paleis docker konteinerius. Priklausomai nuo inteneto ryšio diegimas gali užtrukti nuo 30 min iki kelių valandų.
Sistema bus sudiegta <INSTALL_DIR> direktorijoje

## Patikrinimas

Atidarykite URL naršyklėje: 

Patikrinkite ar visi servisai veikia su docker-compose


