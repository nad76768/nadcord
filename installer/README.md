# Nadcord Installer

Windows app that installs, repairs and uninstalls [Nadcord](https://github.com/nad76768/nadcord).

- `frontend/` — the UI (plain HTML/CSS/JS), rendered with [Wails](https://wails.io) and WebView2
- `app.go` — the methods the UI calls (install, repair, uninstall, OpenAsar, folder picker)
- everything else — the install logic from [Vencord Installer](https://github.com/Vencord/Installer)

It's built by `.github/workflows/installer.yml` and published to the `installer` release.

Licensed under GPL-3.0, based on Vencord Installer © Vendicated and Vencord contributors.
