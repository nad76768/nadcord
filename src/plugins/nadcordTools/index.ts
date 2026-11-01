/*
 * Vencord, a Discord client mod
 * Copyright (c) 2026 Nad
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

// First Nadcord-exclusive plugin.
// Use it as a template: copy this folder, rename it and change the name, description and commands.

import { ApplicationCommandInputType, findOption, OptionalMessageOption, sendBotMessage } from "@api/Commands";
import { definePluginSettings } from "@api/Settings";
import { Devs } from "@utils/constants";
import { sendMessage } from "@utils/discord";
import definePlugin, { OptionType } from "@utils/types";

const settings = definePluginSettings({
    channelUrl: {
        type: OptionType.STRING,
        description: "YouTube channel link shared by the /channel command",
        default: "https://www.youtube.com/@MaxRandom"
    },
    defaultMessage: {
        type: OptionType.STRING,
        description: "Message used by /channel when you don't write one",
        default: "Come check out my videos! 🎬"
    }
});

export default definePlugin({
    name: "NadcordTools",
    description: "Nadcord-exclusive tools: /channel to share your YouTube channel and /nadcord for info about the mod",
    tags: ["Commands"],
    authors: [Devs.Nad],
    settings,

    commands: [
        {
            name: "channel",
            description: "Share your YouTube channel link in chat",
            inputType: ApplicationCommandInputType.BUILT_IN,
            options: [OptionalMessageOption],
            execute(options, { channel }) {
                const message = findOption(options, "message", settings.store.defaultMessage);
                sendMessage(channel.id, {
                    content: `${message} ${settings.store.channelUrl}`.trim()
                });
            }
        },
        {
            name: "nadcord",
            description: "Show info about Nadcord (only visible to you)",
            inputType: ApplicationCommandInputType.BUILT_IN,
            execute(_, { channel }) {
                sendBotMessage(channel.id, {
                    content: "**Nadcord** is active ✅\nNad's personal client mod, based on Vencord (GPL-3.0)."
                });
            }
        }
    ]
});
