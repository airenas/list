## Lentelės ##

**request** 

Saugomos užklausos

| Laukas| Tipas | Paskirtis | 
| ---|-|-|
| ID[*pk*] | string | Transkripcijos ID |
| email| string | Vart. el. paštas |
| date | dateTime | Užklausos laikas |
| file | string | Pradinis failo pavadinimas |


*pk* - raktinis laukas

---
**status**

Saugomas dabartinis transkripcijos statusas

| Laukas| Tipas | Paskirtis |
| ---|-|-|
| ID[*pk*] | string | Transkripcijos ID |
| status | string | Statusas |
| error  | string | Klaida |
| errorCode  | string | Klaidos kodas |
| date   | dateTime | Statuso laikas |

---
**result**

Saugomas galutinis transkripcijos rezultatas

| Laukas| Tipas | Paskirtis |
| ---|-|-|
| ID[*pk*] | string | Transkripcijos ID |
| text | string | Transkripcijos rezultatas |

---
**emailLock**

Lentelė skirta sinchronizuoti el. laiškų siuntimą ir užtikrinti, kad laiškas bus išsiųstas ne daugiau, kaip vieną kartą

| Laukas| Tipas | Paskirtis |
| ---|-|-|
| ID[*pk, i*] | string | Transkripcijos ID |
| key[*pk*] | string | Laiško tipas. Galimos reišmės: Started, Finished |
| status | int | Statusas. Galimos reikšmės: 0 - nepradėta, 1 - siunčiama, 2 - išsiųsta

*i* - indeksuojamas laukas