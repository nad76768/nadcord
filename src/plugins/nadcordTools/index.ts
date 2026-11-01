/*
 * Vencord, a Discord client mod
 * Copyright (c) 2026 Nad
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

// Primo plugin esclusivo di Nadcord.
// Usalo come modello: copia questa cartella, rinominala e cambia nome, descrizione e comandi.

import { ApplicationCommandInputType, findOption, OptionalMessageOption, sendBotMessage } from "@api/Commands";
import { definePluginSettings } from "@api/Settings";
import { Devs } from "@utils/constants";
import { sendMessage } from "@utils/discord";
import definePlugin, { OptionType } from "@utils/types";

const settings = definePluginSettings({
    channelUrl: {
        type: OptionType.STRING,
        description: "Link del canale YouTube condiviso dal comando /canale",
        default: "https://www.youtube.com/@MaxRandom"
    },
    defaultMessage: {
        type: OptionType.STRING,
        description: "Messaggio usato da /canale quando non ne scrivi uno",
        default: "Venite a vedere i miei video! 🎬"
    }
});

export default definePlugin({
    name: "NadcordTools",
    description: "Strumenti esclusivi di Nadcord: /canale per condividere il canale YouTube e /nadcord per le info sulla mod",
    tags: ["Commands"],
    authors: [Devs.Nad],
    settings,

    commands: [
        {
            name: "canale",
            description: "Condividi in chat il link del canale YouTube",
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
            description: "Mostra info su Nadcord (visibile solo a te)",
            inputType: ApplicationCommandInputType.BUILT_IN,
            execute(_, { channel }) {
                sendBotMessage(channel.id, {
                    content: "**Nadcord** è attivo ✅\nClient mod personale di Nad, basata su Vencord (GPL-3.0)."
                });
            }
        }
    ]
});
