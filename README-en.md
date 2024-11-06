🇮🇹 [IT](README.md) | 🇬🇧 EN

# doable-webui

A web ui to interact with todos written and synced with [Doable](https://doable.at/), written in Go (web server/api) and JavaScript (client browser).

⚠️ Work in progress, code and features are subject to changes. The code is written in english but the interface is in italian, it will be translated in the future.

🐛 If you find an issue or want a new feature open an [issue](https://github.com/matteolomba/doable-webui-go/issues) or a [pull request](https://github.com/matteolomba/doable-webui-go/pulls)

## Features

- View todos
- Add new todos (to be implemented)
- Edit a todo (to be implemented)
- Remove a todo (to be implemented)

## Requirements

- Doable sync on Nextcloud (or WebDAV, untested) set up and active.
- The Doable folder, which must be in the same directory as the program, must contain the files that Doable syncs. Automatic sync with Nextcloud in the program may be implemented in the future.

## .env

If you want to change the log level, create a `.env` file in the root of the project with the following content:

```env
LOG_LEVEL=DEBUG # Or INFO, WARN, ERROR, FATAL (same as error), default: WARN
```

## Credits

- [Bootstrap](https://getbootstrap.com/) - Used and included in the project, [MIT](https://github.com/twbs/bootstrap/blob/main/LICENSE) license
- [Rubik Font](https://fonts.google.com/specimen/Rubik) - Used and included in the project, [OFL](https://fonts.google.com/specimen/Rubik/license) license
- [Feather Icons](https://github.com/feathericons/feather) - Used and included in the project, [MIT](https://github.com/feathericons/feather/blob/main/LICENSE) license
