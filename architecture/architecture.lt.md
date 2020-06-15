# Transkribatoriaus IT sprendimo archtektūra, servisai

Sistemą sudarantys komponentai pateikti [AFT-Schema.png](AFT-Schema.png) diagramoje. Kiekvienas komponentas realizuoja atskirą servisą ir yra paruoštas kaip [docker](https://www.docker.com/) konteineris. Konteineriai viešai prieinami [https://hub.docker.com/](https://hub.docker.com/). Diegimui naudojama docker compose technologija. Naudojant [docker swarm](https://docs.docker.com/swarm/overview/) arba [Kubernettes](https://kubernetes.io/) vartotojas gali diegti sistemą ant kelių fizinių mašinų. Vienu metu sistemoje gali būti daug veikiančių tų pačių transkripcijos servisų (diagramoje TS: ...). Failų saugyklos (FS: ...) sukonfigūruojamos kaip [Docker volume](https://docs.docker.com/storage/volumes/) arba [Kubernettes PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/).

## Servisų aprašymas

| Servisas/programa| Paskirtis | Realizuojama komponentais | Docker komponento versija |
| ---|-|-|-|
| Failų priėmimo, rezultatų pateikimo programa | Naršyklei skirtas puslapis, nusiųsti failą į sistemą, gauti transkripcijos statusą ir rezultatus  | [Angular](https://angular.io/), nginx ||
| Semantika integracijos servisas | Integruoja semantikos paslaugų puslapį su Transkripcijos IT sprendimu. Klauso Kafka pranešimų eilės, saugo transkripcijos rezultatus Fonogramų saugykloje | Go ||
| Failų priėmimo servisas | Realizuoja sąsają priimti failus iš vartotojo. [API](https://app.swaggerhub.com/apis/aireno/Transkipcija/1.4.0) | Go |
| Transkribavimo statuso peržiūros servisas | Grąžina informaciją apie transkribuojamą audio failą. [API](https://app.swaggerhub.com/apis/aireno/Transkipcija/1.4.0) | Go ||
| Transkribavimo rezultatų pateikimo servisas | Pateikia informaciją apie transkribuotą audio failą: failas, transkripcija, kitos transkripcijos hipotezės. [API](https://app.swaggerhub.com/apis/aireno/Transkipcija/1.4.0) | Go ||
| DB servisas | Saugo duomenis apie audio failo transkribavimo statusą, rezultatus, klaidas | [Mongo DB](https://www.mongodb.com) | mongo:4.1.1 |
| Įvykių eilės servisas | Laiko įvykių užklausas, persiunčia servisams užduotis | [RabbitMQ](https://www.rabbitmq.com/) | rabbitmq:3.7-management |
| Transkribavimo valdymo servisas | Realizuoja transkripcijos logiką: žingsnius ką reikia padaryti nuo audio failo gavimo iki vartotojo informavimo apie rezultatus/klaidas. Skirsto darbus kitiems sistemos servisams per įvykių eilę, saugo transkripcijos statusą į DB. | Go ||
| TS: Pradinio audio failų apdorojimo servisas | Transformuoja vartotojo įkeltą failą į vieningą audio formatą reikalingą transkripcijai | Go, [ffmep](https://www.ffmpeg.org/), [sox](http://sox.sourceforge.net/), Bash |
| TS: Audio failo skaidymo pagal kalbėtojo tipą servisas | Suranda ir suskaido audio failą pagal kalbėtojo/kalbėjimo tipus: tyla, moteriškas, vyriškas balsas ir pan. | Go, [LIUM_SpkDiarization-4.2](https://projets-lium.univ-lemans.fr/spkdiarization/download/), Bash |
| TS: Transkribavimo servisas | Skaičiuoja parametrus ir dekoduoja audio signalą, naudojamas Kaldi paketas. Apskaičiuoja kelias geriausias hipotezes.| Go, [Kaldi](http://kaldi-asr.org/), Bash |
| TS: Geriausios hipotezės parinkimo servisas | Naudojant RNNLM kalbos modelį parenkamas geriausias transkripcijos variantas | Go, Kaldi, Bash |
| TS: Rezultato paruošimo servisas | Tekste sudėlioja skyrybos ženklus, paruošia grafą redagavimo, subtitrų formatais | Go, Kaldi, Bash, Perl ||
| Darbų paskirstymo servisas | Skirsto transkribavimo ir geriausios hipotezės atrinkimo užduotis pagal vartotojo nurodytą atpažintuvo parametrą. Optimizuoja bendrą sistemos greitaveiką | Go ||
| Skyrybos ženklų atstatymo servisas | Tekste sudėlioja skyrybos ženklus | Go ||
| Tensorflow servisas | Teikia paslaugas tensorflow skyrybos ženklų atstatymo medeliui | [Tensorflow serving](https://www.tensorflow.org/tfx/guide/serving) | tensorflow/serving:1.14.0
| Vartotojo informavimo servisas | Siunčia el.laišką apie transkripcijos startą, pabaigą vartotojui | Go |
| Nereikalingų vartotojo duomenų trynimo servisas | Ištrina transkribavimo duomenis iš failinės sistemos, kai jie jau nereikalingi. | Go
| Metrikų surinkimo servisas | Realizuoja sąsają metrikų registravimui iš programinių skriptų | Go  ||
| Sistemos darbo analizės servisas | Bus kaupiama sistemos darbo statistika: užklausų kiekis, apdorojimo laikai, servisų naudojami resursai (atmintis, procesoriaus laikas) | [Prometheus](https://prometheus.io/) | prom/prometheus:v2.17.2, prom/node-exporter:v0.18.1
| FS: DB duomenų saugykla | Saugomi DB ir Įvykių eilės servisų failai | [Docker volume](https://docs.docker.com/storage/volumes/) arba [Kubernettes PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
| FS: Transkribavimo failų saugykla | Saugomi visi darbiniai failai nuo pirmo vartotojo pateikto audio failo iki galutinių transkripcijos tekstų ir hipotezių failų. | [Docker volume](https://docs.docker.com/storage/volumes/) arba [Kubernettes PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) |
| FS: Atpažinimo modulių failų saugykla | Failai, kuriuose saugomi akustiniai, kalbos, skyrybos ženklų atstatymo modeliai, MFCC skaičiavimo ir kiti konfigūraciniai failai | [Docker volume](https://docs.docker.com/storage/volumes/) arba [Kubernettes PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
