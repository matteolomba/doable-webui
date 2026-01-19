🇮🇹 IT | 🇬🇧 [EN](README-en.md)

# doable-webui

Un'interfaccia web per interagire con cose da fare scritte e sincronizzate con [Doable](https://doable.at/), scritta in Go (web server/api) e JavaScript.

⚠️ Sviluppo in corso, codice e funzioni soggetti a cambiamenti. Il codice è scritto in inglese ma l'interfaccia è in italiano, verrà tradotta anche in inglese un giorno.

🐛 Se trovi un problema o vorresti una nuova funzionalità apri un [issue](https://github.com/matteolomba/doable-webui-go/issues) o una [pull request](https://github.com/matteolomba/doable-webui-go/pulls)

## Cosa ci puoi fare

- ⚡ **Gestione Cose da fare (CRUD)**: Visualizza, aggiungi, modifica, completa ed elimina le cose da fare
- 📋 **Gestione Liste (CRUD)**: Crea, modifica ed elimina liste
- 🎨 **Personalizzazione liste**: Icone per le liste (visualizzate solo sulla webui)
- 🔍 **Filtri**: Per lista, stato e ricerca testuale
- 🌙 **Tema chiaro/scuro**: Attivabile/disattivabile manualmente

## Requisiti

- Sincronizzazione di Doable con Nextcloud (o WebDAV, non testato) impostata e attiva.
- La cartella sync, che deve essere presente nella stessa directory del programma, deve contenere i file che Doable sincronizza. In futuro potrebbe venire implementata la sincronizzazione automatica con Nextcloud direttamente nel programma.

## .env

Se vuoi cambiare il livello dei log, crea un file `.env` nella root del progetto con il seguente contenuto:

```env
LOG_LEVEL=DEBUG # Oppure INFO, WARN, ERROR, FATAL (uguale a error), default: WARN
```

## Crediti

- [Tailwind CSS](https://tailwindcss.com/) - Utilizzato nel progetto, licenza [MIT](https://github.com/tailwindlabs/tailwindcss/blob/master/LICENSE)
- [Alpine.js](https://alpinejs.dev/) - Utilizzato nel progetto, licenza [MIT](https://github.com/alpinejs/alpine/blob/main/LICENSE.md)
- [Outfit Font](https://fonts.google.com/specimen/Outfit) - Utilizzato e incluso nel progetto, licenza [OFL](https://fonts.google.com/specimen/Outfit/license)
- [Feather Icons](https://github.com/feathericons/feather) - Utilizzato e incluso nel progetto (SVG inline direttamente nel codice), licenza [MIT](https://github.com/feathericons/feather/blob/main/LICENSE)
