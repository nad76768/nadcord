/*
 * MaxCord, a fork of Vencord - a modification for Discord's desktop app
 * Copyright (c) 2026 Max Random
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

// Primo plugin esclusivo di MaxCord.
// Usalo come modello: copia questa cartella, rinominala e cambia nome, descrizione e comandi.

import { ApplicationCommandInputType, findOption, OptionalMessageOption, sendBotMessage } from "@api/Commands";
import { definePluginSettings } from "@api/Settings";
import { sendMessage } from "@utils/discord";
import definePlugin, { OptionType } from "@utils/types";

const MaxRandom = { name: "Max Random", id: 0n };

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
    name: "MaxCordTools",
    description: "Strumenti esclusivi di MaxCord: /canale per condividere il canale YouTube e /maxcord per le info sulla mod",
    tags: ["Commands"],
    authors: [MaxRandom],
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
            name: "maxcord",
            description: "Mostra info su MaxCord (visibile solo a te)",
            inputType: ApplicationCommandInputType.BUILT_IN,
            execute(_, { channel }) {
                sendBotMessage(channel.id, {
                    content: "**MaxCord** è attivo ✅\nClient mod personale di Max Random, basata su Vencord (GPL-3.0)."
                });
            }
        }
    ]
});
