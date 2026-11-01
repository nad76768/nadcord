@echo off
chcp 65001 >nul
title Installazione Nadcord
cd /d "%~dp0"

echo === Nadcord - installazione ===
echo.

where node >nul 2>nul
if errorlevel 1 (
    echo [ERRORE] Node.js non trovato. Scaricalo da https://nodejs.org ^(versione 22 o superiore^) e riprova.
    pause
    exit /b 1
)

where git >nul 2>nul
if errorlevel 1 (
    echo [ERRORE] Git non trovato. Scaricalo da https://git-scm.com e riprova.
    pause
    exit /b 1
)

where pnpm >nul 2>nul
if errorlevel 1 (
    echo Installo pnpm...
    call npm install -g pnpm
    if errorlevel 1 goto errore
)

echo [1/3] Installo le dipendenze...
call pnpm install --frozen-lockfile
if errorlevel 1 goto errore

echo [2/3] Compilo Nadcord...
call pnpm build
if errorlevel 1 goto errore

echo [3/3] Avvio l'installer: scegli la tua installazione di Discord.
call pnpm inject
if errorlevel 1 goto errore

echo.
echo Fatto! Riapri Discord e cerca la sezione Nadcord nelle impostazioni.
pause
exit /b 0

:errore
echo.
echo [ERRORE] Qualcosa e' andato storto. Copia il testo qui sopra e mandamelo.
pause
exit /b 1
