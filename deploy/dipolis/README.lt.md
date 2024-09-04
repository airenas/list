# Diegimas vienoje mašinoje naudojant *Docker*

DiPolis versija

## Apie

Transkribatroriaus IT sprendimas yra realizuotas *Docker* komponentais. Visa sistema sukonfigūruota ir paruošta paleisti su *docker compose* konfigūraciniu failu. Sistemos darbui taip pat reikalingi kai kurie papildomi (nedokerizuoti) binariniai vykdomieji failai ir lietuvių kalbos atpažinimo modelių failai. Diegiant Jums reikės:

- atlikti pakeitimus konfiguraciniame faile,

- nauojantis paruoštų skriptų pagalba, parsiųsti reikalingus binarinius ir modelių failus,
- su *docker ompose* startuoti sistemą.
Transcribatroriaus IT sprendimas gali būti diegiamas, bet kurioje operacinėje sistemoje, kuri palaiko Docker technologiją, bet buvo testuotas ir šis aprašymas apima tik *Linux* sistemas. *Win* ir *Mac* operacinėse sistemose tikėtina bus reikalingas papildomas *docker-compose.yml* failo pritaikymas.

## Reikalavimai

Aparatūrai:

| Komponen-tas | Min reikalavimai | Rekomenduo-jama | Papildomai |
| -----------------|------------------|---------------------|-------------------------------------------|
| Platform | x86_64 | | |
| CPU | 64-bit, 2 branduoliai | 8 branduoliai | |
| HDD | 40 Gb | | Priklausomai nuo sudiegtų atpažinimo modelių. Vienam modeliui papildomai reikia apie 10 Gb |
| RAM | 32 Gb | 48 Gb | |

Operacinė sistema: Linux OS 64-bit (papildomai žiūrėkite [reikalavimus Docker instaliacijai](https://docs.docker.com/engine/install/)). Turi būti sudiegta:

| Komponentas | Min versija | URL |
| ---|-|-|
| Docker | 27.1.2 | [Link](https://docs.docker.com/engine/install/)

Papildomi įrankiai naudojami instaliuojant: [make](https://www.gnu.org/software/make/manual/make.html), [git](https://git-scm.com/download/linux), [wget](https://www.gnu.org/software/wget/manual/wget.html), [tar](https://www.gnu.org/software/tar/manual/).

## Prieš diegiant

Patikrinkite ar visi reikalingi komponentai veikia mašinoje:

```bash
    ## Docker
    docker run hello-world
    ## Kiti komponentai
    make --version
    tar --version
    wget --version
    git --version
```

## Diegimas

1. Parsisiųskite diegimo skriptus (ši git repositorija):

    `git clone https://github.com/airenas/list.git`

    `cd list/deploy/dipolis`

    Docker diegimo skriptai yra direktorijoje yra *list/deploy/dipolis*.

1. Pasirinkite diegimo versiją:

    `git checkout xxx`   
    
    pateiks diegėjas

1. Paruoškite konfigūracinį diegimo failą *Makefile.options*:

    `cp Makefile.options.template Makefile.options`

1. Sukonfigūruokite *Makefile.options*:

    | Parametras | Priva-lomas | Paskirtis | Pvz |
    |------------------|-----|-----------------------------------|------------------|
    | *deploy_dir* | + | Pilnas kelias iki instaliavimo direktorijos mašinoje. Šioje direktorijoje bus atsiųsti modeliai, sukurtas pakatalogis darbiniams transkribatoriaus failams | /home/user/list
    | *models* | + | Instaliuojami modeliai. Galimi pasirinkimai: *ben-tel-2.0r1* | ben-tel-2.0r1 |
    | rabbitmq_pass | + | Eilės serviso slaptažodis. Nurodykite slaptažodį, kurį servisai naudos prisijungimui prie eilės serviso. Pvz.: sugeneruokite su `pwgen 20 1` ||
    | mongo_pass | + | DB slaptažodis. Nurodykite slaptažodį, kurį servisai naudos prisijungimui prie vidinės DB. Pvz.: sugeneruokite su `pwgen 20 1` ||
    | http_port | + | HTTP portas, kuriuo bus pasiekiami servisai mašinoje | 80 |
    | host_external_url | - | Kompiuterio URL, kuriuo servisai pasiekiami iš išorės. Naudojama nuorodai el. laiške | <https://airenas.eu:7054> |
    | smtp_host | - |SMTP serveris, laiškų siuntimui | 80 |
    | smtp_port | - |SMTP portas | 587 |
    | smtp_username | - | SMTP serverio vartotojas. Jei tuščias - sistema nesiųs informacinių laiškų | olia@gmail.com |
    | smtp_password | - | SMTP slaptažodis |  |
    | smtp_type     | - | SMTP serverio autentifikavimo tipas. Galimos reikšmės: NO_AUTH (kai SMTP serveris nereikalauja slaptažodžio), PLAIN_AUTH (veikia daugumai SMTP serverių, naudoja TLS, jei serveris palaiko), LOGIN (kai SMTP serveris reikalauja Login autentifikacijos) | PLAIN_AUTH |
    | hf_api_token | | jei naudojamas pyannote diarizatorius - https://huggingface.co/pyannote/speaker-diarization |
    | *share_pass* | + | Modelių parsisiuntimo slaptažodis (pateiks VDU) | xxx123 |

1. *optional* Jei norite naudoti pyannote diarizatorių:

    **Nerekomenduojama jei sistemoje nėra GPU.**

    - Sukonfigūruokite `hf_api_token` - žr. https://huggingface.co/pyannote/speaker-diarization, modelis 2.1.
    - Užkomentuokite `diarization-service`, faile `docker-compose.yml`
    - Atkomentuokite `diarization-pyannote-service`, faile `docker-compose.yml`

1. Instaliuokite

    `make install -j4`

    Skriptas parsiųs reikalingus failus, paleis *docker* konteinerius. Priklausomai nuo interneto ryšio diegimas gali užtrukti nuo 30 min iki kelių valandų.
    Sistema bus sudiegta *<deploy_dir>* direktorijoje

## Patikrinimas

1. Patikrinkite ar visi servisai veikia su *docker compose*: `make status`. Visi servisai turi būti *Up* būsenoje.

1. Patikrinkite ar servisas gali priimti užklausas: `make status-service`. Turi grąžinti užklausos kodą 200.

1. Atidarykite URL naršyklėje: *<host_external_url>/ausis/*. Turi atsidaryti demo puslapis.

## Servisų sustabdymas/valdymas

Servisai valdomi su *docker compose* komanda:

```bash
    cd <deploy_dir>
    ##Sustabdymas
    docker compose stop
    ##Paleidimas
    docker compose up -d
```

## Duomenų atnaujinimas

1. Atnaujinus duomenis, bus pakeista ir ši repositorija su nuorodomis į naujus modelių failus. Patikrinkite, kad turite naujausius skriptus:

    `git pull`

1. Pasirinkite norimą versiją:

    `git checkout <VERSIJA>`

    Versija turi priskirtą *git* žymą. Galimas versijas galite sužinoti su komanda: `git tag`.

1. Atnaujinkite servisus/duomenis:

```bash
    make clean-docker
    make install -j4
```

## Pašalinimas

```bash
    make clean
```

Komandą reikia vykdyti *admin* teisėmis, pvz.: `sudo sh -c 'make clean'`.
