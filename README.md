# Nadcord

Nad's personal Discord client mod, based on [Vencord](https://github.com/Vendicated/Vencord).

Includes every Vencord plugin and feature (150+ plugins, themes, QuickCSS, updater, settings sync), plus Nadcord-exclusive plugins.

## Exclusive plugins

| Plugin | What it does |
| --- | --- |
| `NadcordTools` | `/channel` shares your YouTube channel link, `/nadcord` shows info about the mod |

Exclusive plugins live in `src/plugins/` alongside the Vencord ones. To create a new one, copy `src/plugins/nadcordTools`, rename it and edit it.

## Installing (Windows)

1. Download **[NadcordInstaller.exe](https://github.com/nad76768/nadcord/releases/download/installer/NadcordInstaller.exe)**.
2. Fully close Discord (including from the system tray).
3. Open the installer, pick your Discord and press **Install**.
4. Reopen Discord: you'll find the **Nadcord** section in User Settings.

The same installer can also **Reinstall / Repair** or **Uninstall** Nadcord.

> Windows may show "Windows protected your PC" because the installer is not code-signed. Click **More info** → **Run anyway**.

The installer's source code is in the `installer/` folder.

### Developer install

Requirements: [Node.js](https://nodejs.org) 22+ and [Git](https://git-scm.com).
Double-click `install-nadcord.bat` to build and inject Nadcord from your local source.
To uninstall: `pnpm uninject`.

## Updates

- **Nadcord:** every push to `main` is built automatically and published to the `devbuild` release. Use the built-in updater or **Reinstall / Repair** in the installer.
- **Installer:** changes in `installer/` are built automatically and published to the `installer` release (same download link).
- **Upstream Vencord:** `git pull upstream main`, then push.

## License and credits

Nadcord is a fork of Vencord © Vendicated and contributors, licensed under **GPL-3.0-or-later** (see `LICENSE`).
If you distribute Nadcord, its source code must remain public under the same license.
The original Vencord README is in `README.vencord.md`.

> Client mods violate Discord's Terms of Service. Use at your own risk.
