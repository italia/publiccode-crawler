<!--
    This file is linked from https://developers.italia.it,
    DON'T REMOVE/RENAME it without updating the link first.
-->

## Indice di vitalità

Questo documento illustra le modalità di calcolo dell'`indice di vitalità` che
viene visualizzato per ogni software indicizzato nel [catalogo del
riuso](https://developers.italia.it/it/software).

Questo indice rappresenta la vitalità di un certo repository nell'ultimo
frangente di tempo.
Infatti mette a fuoco i seguenti parametri:

* Code activity: il numero di commit e merge effettuati giornalmente;
* Release history: il numero di rilasci giornalieri;
* User community: il numero di autori unici;
* Longevity: l'età del progetto.

Per ognuno di questi parametri vengono assegnati dei punteggi che poi saranno
sommati per formare l'indicatore finale (`vitality index`).
In questo momento l'algoritmo usa una finestra di `60 giorni` per il calcolo,
ovvero, per ognuna delle succitate categorie, vengono prese in considerazione
le azioni effettuate negli ultimi due mesi.

### Code Activity

Questo indicatore rappresenta il numero di attività effettuate all'interno
della finestra temporale analizzata. In particolare, due azioni sono
considerate:

1. il numero di commit;

2. il numero di merge.

Questo indicatore ha un valore minimo di 2, quando vi sono state da
0 a 4 azioni, e un massimo di 60, quando le attività sono più di 35.
Per avere ulteriori informazioni è possibile visitare
[questo](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/vitality-ranges.yml#L31-L62)
file.

### Release History

Questo parametro rappresenta il numero di rilasci effettuati nell'ultimo
frangente di tempo.
Analizzando un repository git un modo di conoscere il numero di rilasci
è quello di analizzare i `tag`. Questo parametro varia da un valore minimo di 20,
punti quando vi sono stati da 0 a 1 rilasci, ad un massimo di 50 per un numero
di rilasci che varia da 4 a 100.
Per ulteriori informazioni si può visitare [questo](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/vitality-ranges.yml#L64-L77)
file.

### User community

Questo indicatore rappresenta il numero di utenti che ha eseguito un'operazione
di `commit` sul repository nell'arco di tempo analizzato.
Sapere quanti utenti hanno effettivamente lavorato sul codice fornisce
un'indicazione della community che si è formata intorno ad un dato prodotto.
Questo parametro ha un valore minimo di 4 e un massimo di 36.
Maggiori informazioni sono reperibili [qui](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/vitality-ranges.yml#L1-L29).

### Longevity

La longevità è l'età del repository. Questo parametro varia dunque in base
a quando il repository sia stato aperto.
Il valore di questo indicatore varia da 20 punti ad un massimo di 35 punti.
[Qui](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/vitality-ranges.yml#L79-L89)
è possibile reperire maggiori informazioni in merito al range.

### Calcolo finale

Tutti i 4 precedenti valori non sono valori assoluti ma sono calcolati su base
giornaliera, dunque alla fine della rilevazione vi saranno **n** valori dove
**n** è il numero dei giorni presi in considerazione. Per ottenere il valore
finale si effettua quindi una media di ogni indicatore e una somma finale che
costituisce l'indice che viene visualizzato nel catalogo.
[Vedi maggiori dettagli](https://github.com/italia/publiccode-crawler/blob/663c661ca3b0d6e1578f24c7be97fd35e28abe87/crawler/crawler/repo_activity.go#L100-L117)
