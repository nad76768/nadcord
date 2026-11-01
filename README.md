# Nadcord

Client mod personale per Discord di **Nad**, basata su [Vencord](https://github.com/Vendicated/Vencord).

Include tutti i plugin e le funzioni di Vencord (oltre 150 plugin, temi, QuickCSS, updater, sincronizzazione impostazioni), più i plugin esclusivi di Nadcord.

## Plugin esclusivi

| Plugin | Cosa fa |
| --- | --- |
| `NadcordTools` | `/canale` condivide il link del canale YouTube, `/nadcord` mostra le info sulla mod |

I plugin esclusivi stanno in `src/plugins/` come quelli di Vencord. Per crearne uno nuovo, copia la cartella `src/plugins/nadcordTools`, rinominala e modificala.

## Installazione (Windows)

1. Scarica **[NadcordInstaller.exe](https://github.com/nad76768/nadcord/releases/download/installer/NadcordInstaller.exe)**.
2. Chiudi Discord del tutto (anche dall'icona vicino all'orologio).
3. Apri l'installer, scegli il tuo Discord e premi **Installa**.
4. Riapri Discord: nelle impostazioni trovi la sezione **Nadcord**.

Dallo stesso installer puoi anche **Reinstallare / Riparare** o **Disinstallare** Nadcord.

Il codice dell'installer si trova nella cartella `installer/`.

### Installazione da sviluppatore

Requisiti: [Node.js](https://nodejs.org) 22+ e [Git](https://git-scm.com). Fai doppio clic su `installa-nadcord.bat` (compila dal codice locale).
Per disinstallare: `pnpm uninject`.

## Aggiornamenti

- **Dal tuo repository:** l'updater integrato controlla il repository GitHub da cui hai compilato.
- **Dal Vencord originale:** `git pull upstream main`, poi ricompila con `pnpm build`.

## Licenza e crediti

Nadcord è un fork di Vencord © Vendicated e collaboratori, distribuito sotto licenza **GPL-3.0-or-later** (vedi `LICENSE`).
Se distribuisci Nadcord ad altri, il codice sorgente deve restare pubblico con la stessa licenza.
Il README originale di Vencord si trova in `README.vencord.md`.

> Le client mod violano i Termini di Servizio di Discord. Usala consapevolmente.
