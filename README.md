# MaxCord

Client mod personale per Discord di **Max Random**, basata su [Vencord](https://github.com/Vendicated/Vencord).

Include tutti i plugin e le funzioni di Vencord (oltre 150 plugin, temi, QuickCSS, updater, sincronizzazione impostazioni), più i plugin esclusivi di MaxCord.

## Plugin esclusivi

| Plugin | Cosa fa |
| --- | --- |
| `MaxCordTools` | `/canale` condivide il link del canale YouTube, `/maxcord` mostra le info sulla mod |

I plugin esclusivi stanno in `src/plugins/` come quelli di Vencord. Per crearne uno nuovo, copia la cartella `src/plugins/maxCordTools`, rinominala e modificala.

## Installazione (Windows)

Requisiti: [Node.js](https://nodejs.org) versione 22 o superiore e [Git](https://git-scm.com).

1. Chiudi Discord.
2. Fai doppio clic su `installa-maxcord.bat`.
3. Nella finestra dell'installer scegli la tua installazione di Discord.
4. Riapri Discord: nelle impostazioni trovi la sezione **MaxCord**.

Per disinstallare: `pnpm uninject`.

## Aggiornamenti

- **Dal tuo repository:** l'updater integrato controlla il repository GitHub da cui hai compilato.
- **Dal Vencord originale:** `git pull upstream main`, poi ricompila con `pnpm build`.

## Licenza e crediti

MaxCord è un fork di Vencord © Vendicated e collaboratori, distribuito sotto licenza **GPL-3.0-or-later** (vedi `LICENSE`).
Se distribuisci MaxCord ad altri, il codice sorgente deve restare pubblico con la stessa licenza.
Il README originale di Vencord si trova in `README.vencord.md`.

> Le client mod violano i Termini di Servizio di Discord. Usala consapevolmente.
