## Lentelės ##

**request** 
Saugomos užklausos

| Laukas| Tipas | Paskirtis | 
| ---|-|-|
| ID<sup>pk</sup> | string | Transkripcijos ID |
| email| string | Vart. el. paštas |
| date | dateTime | Užklausos laikas |

<sup>pk</sup> - raktinis laukas

---
**status**
Saugomas dabartinis transkripcijos statusas

| Laukas| Tipas | Paskirtis |
| ---|-|-|
| ID<sup>pk</sup> | string | Transkripcijos ID |
| status | string | Statusas |
| error  | string | Klaida |
| date   | dateTime | Statuso laikas |

---
**result**
Saugomas galutinis transkripcijos rezultatas

| Laukas| Tipas | Paskirtis |
| ---|-|-|
| ID<sup>pk</sup> | string | Transkripcijos ID |
| text | string | Transkripcijos rezultatas |

---
**emailLock**
Lentelė skirta sinchronizuoti el. laiškų siuntimą ir užtikrinti, kad laiškas bus išsiųstas ne daugiau, kaip vieną kartą

| Laukas| Tipas | Paskirtis |
| ---|-|-|
| ID<sup>pk, i</sup> | string | Transkripcijos ID |
| key<sup>pk</sup> | string | Laiško tipas. Galimos reišmės: Started, Finished |
| status | int | Statusas. Galimos reikšmės: 0 - nepradėta, 1 - siunčiama, 2 - išsiųsta
<sup>i</sup> - indeksuojamas laukas