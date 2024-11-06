🇮🇹 IT | 🇬🇧 [EN](README-en.md)

# doable-webui

Un'interfaccia web per interagire con cose da fare scritte e sincronizzate con [Doable](https://doable.at/), scritta in Go (web server/api) e JavaScript.

⚠️ Sviluppo in corso, codice e funzioni soggetti a cambiamenti. Il codice è scritto in inglese ma l'interfaccia è in italiano, verrà tradotta anche in inglese un giorno.

🐛 Se trovi un problema o vorresti una nuova funzionalità apri un [issue](https://github.com/matteolomba/doable-webui-go/issues) o una [pull request](https://github.com/matteolomba/doable-webui-go/pulls)

## Cosa ci puoi fare

- Visualizzare le cose da fare
- Aggiungere nuove cose da fare (da implementare)
- Modificare una cosa da fare (da implementare)
- Rimuovere una cose da fare (da implementare)

## Requisiti

- Sincronizzazione di Doable con Nextcloud (o WebDAV, non testato) impostata e attiva.
- La cartella Doable, che deve essere presente nella stessa directory del programma, deve contenere i file che Doable sincronizza. In futuro potrebbe venire implementata la sincronizzazione automatica con Nextcloud direttamente nel programma.

## .env

Se vuoi cambiare il livello dei log, crea un file `.env` nella root del progetto con il seguente contenuto:

```env
LOG_LEVEL=DEBUG # Oppure INFO, WARN, ERROR, FATAL (uguale a error), default: WARN
```

## Crediti

- [Bootstrap](https://getbootstrap.com/) - Utilizzato e incluso nel progetto, licenza [MIT](https://github.com/twbs/bootstrap/blob/main/LICENSE)
- [Rubik Font](https://fonts.google.com/specimen/Rubik) - Utilizzato e incluso nel progetto, licenza [OFL](https://fonts.google.com/specimen/Rubik/license)
- [Feather Icons](https://github.com/feathericons/feather) - Utilizzato e incluso nel progetto, licenza [MIT](https://github.com/feathericons/feather/blob/main/LICENSE)
