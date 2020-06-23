# *Docker* konteinerių paruošimas

Transkribatoriaus IT sprendimas susideda iš daugelio atskirų servisų, kurie yra konteinerizuoti *docker* technologijos pagalba.
Kodo kompiliavimas ir *docker* konteinerių paruošimas realizuotas *make* skriptų pagalba. Kiekvienam servisui paruoštas atskiras skriptas direktorijoje *deploy/local/<serviso_vardas>*

## Prieš

1. Parsisiųskite ir paruoškite servisų kodo repozitoriją [bitbucket.org/airenas/listgo](https://bitbucket.org/airenas/listgo). Instrukcija [čia](https://bitbucket.org/airenas/listgo/src/master/README.lt.md).

1. Sukonfigūruokite [deploy/local/Makefile.options](Makefile.options). Nustatykite direktorijas:

    - *GO_SRC_DIR* - lokali direktorija, kur parsiųsta [bitbucket.org/airenas/listgo](https://bitbucket.org/airenas/listgo)
    - *SRC_DIR* - nuoroda į šios repozitorijos *src* direktoriją

## Konteinerio paruošimas

1. Pakeiskite serviso versiją, jei reikia, faile [deploy/local/Makefile.options](Makefile.options)

1. Paruoškite konteinerį:

    `cd deploy/local/<serviso_vardas>`

    `make clean dbuild`

1. Arba išsaugokite konteinerį [Docker Hub](https://hub.docker.com/) (prieš tai turite buti prisijungti prie *Docker Hub* su `docker login` komanda):

    `cd deploy/local/<serviso_vardas>`

    `make clean dpush`
