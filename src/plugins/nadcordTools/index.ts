/*
 * Nadcord, a fork of Vencord - a modification for Discord's desktop app
 * Copyright (c) 2026 Nad
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Primo plugin esclusivo di Nadcord.
// Usalo come modello: copia questa cartella, rinominala e cambia nome, descrizione e comandi.

import { ApplicationCommandInputType, findOption, OptionalMessageOption, sendBotMessage } from "@api/Commands";
import { definePluginSettings } from "@api/Settings";
import { sendMessage } from "@utils/discord";
import definePlugin, { OptionType } from "@utils/types";

const Nad = { name: "Nad", id: 0n };

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
    authors: [Nad],
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
