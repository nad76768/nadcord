@echo off
chcp 65001 >nul
title Nadcord developer install
cd /d "%~dp0"

echo === Nadcord - developer install ===
echo.

where node >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Node.js not found. Download it from https://nodejs.org ^(version 22 or newer^) and try again.
    pause
    exit /b 1
)

where git >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Git not found. Download it from https://git-scm.com and try again.
    pause
    exit /b 1
)

where pnpm >nul 2>nul
if errorlevel 1 (
    echo Installing pnpm...
    call npm install -g pnpm
    if errorlevel 1 goto failed
)

echo [1/3] Installing dependencies...
call pnpm install --frozen-lockfile
if errorlevel 1 goto failed

echo [2/3] Building Nadcord...
call pnpm build
if errorlevel 1 goto failed

echo [3/3] Starting the installer: pick your Discord install.
call pnpm inject
if errorlevel 1 goto failed

echo.
echo Done! Reopen Discord and look for the Nadcord section in User Settings.
pause
exit /b 0

:failed
echo.
echo [ERROR] Something went wrong. Check the output above.
pause
exit /b 1
