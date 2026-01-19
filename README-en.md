🇮🇹 [IT](README.md) | 🇬🇧 EN

# doable-webui

A web ui to interact with todos written and synced with [Doable](https://doable.at/), written in Go (web server/api) and JavaScript (client browser).

⚠️ Work in progress, code and features are subject to changes. The code is written in english but the interface is in italian, it will be translated in the future.

🐛 If you find an issue or want a new feature open an [issue](https://github.com/matteolomba/doable-webui-go/issues) or a [pull request](https://github.com/matteolomba/doable-webui-go/pulls)

## Features

- ⚡ **Todo Management (CRUD)**: View, add, edit, complete and delete todos
- 📋 **List Management (CRUD)**: Create, edit and delete lists
- 🎨 **List Customization**: Icons for lists (displayed only on the webui)
- 🔍 **Filters**: By list, status and text search
- 🌙 **Light/Dark Theme**: Manually togglable

## Requirements

- Doable sync on Nextcloud (or WebDAV, untested) set up and active.
- The sync folder, which must be in the same directory as the program, must contain the files that Doable syncs. Automatic sync with Nextcloud in the program may be implemented in the future.

## .env

If you want to change the log level, create a `.env` file in the root of the project with the following content:

```env
LOG_LEVEL=DEBUG # Or INFO, WARN, ERROR, FATAL (same as error), default: WARN
```

## Credits

- [Tailwind CSS](https://tailwindcss.com/) - Used in the project, [MIT](https://github.com/tailwindlabs/tailwindcss/blob/master/LICENSE) license
- [Alpine.js](https://alpinejs.dev/) - Used in the project, [MIT](https://github.com/alpinejs/alpine/blob/main/LICENSE.md) license
- [Outfit Font](https://fonts.google.com/specimen/Outfit) - Used and included in the project, [OFL](https://fonts.google.com/specimen/Outfit/license) license
- [Feather Icons](https://github.com/feathericons/feather) - Used and included in the project (inline SVGs directly in the code), [MIT](https://github.com/feathericons/feather/blob/main/LICENSE) license
