/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"vencordinstaller/buildinfo"
)

const ReleaseUrl = "https://api.github.com/repos/nad76768/nadcord/releases/tags/devbuild"
const ReleaseUrlFallback = "https://api.github.com/repos/nad76768/nadcord/releases/tags/devbuild"

var UserAgent = "NadcordInstaller/" + buildinfo.InstallerGitHash + " (https://github.com/nad76768/nadcord)"

var LinuxDiscordNames = []string{
	"Discord",
	"DiscordPTB",
	"DiscordCanary",
	"DiscordDevelopment",
	"discord",
	"discordptb",
	"discordcanary",
	"discorddevelopment",
	"discord-ptb",
	"discord-canary",
	"discord-development",
	// Flatpak
	"com.discordapp.Discord",
	"com.discordapp.DiscordPTB",
	"com.discordapp.DiscordCanary",
	"com.discordapp.DiscordDevelopment",
}
