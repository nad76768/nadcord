/*
 * Vencord, a Discord client mod
 * Copyright (c) 2025 Vendicated and contributors
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

import { useLayoutEffect } from "@webpack/common";

import { useForceUpdater } from "./react";

const cssRelevantDirectives = ["style-src", "style-src-elem", "img-src", "font-src"] as const;

export const CspBlockedUrls = new Set<string>();

// Payment and anti-fraud trackers that Discord itself loads on its billing pages.
// They are blocked on purpose and have nothing to do with themes, so don't list them.
const IGNORED_HOSTS = ["paypal.com", "paypalobjects.com", "braintreegateway.com", "braintree-api.com", "stripe.com", "stripe.network", "adyen.com"];

function isDiscordInternalResource(url: string) {
    try {
        const { hostname } = new URL(url);
        return IGNORED_HOSTS.some(h => hostname === h || hostname.endsWith("." + h));
    } catch {
        return false;
    }
}

const CspErrorListeners = new Set<() => void>();

document.addEventListener("securitypolicyviolation", ({ effectiveDirective, blockedURI }) => {
    if (!blockedURI || !cssRelevantDirectives.includes(effectiveDirective as any)) return;
    if (isDiscordInternalResource(blockedURI)) return;

    CspBlockedUrls.add(blockedURI);

    CspErrorListeners.forEach(listener => listener());
});

export function useCspErrors() {
    const forceUpdate = useForceUpdater();

    useLayoutEffect(() => {
        CspErrorListeners.add(forceUpdate);

        return () => void CspErrorListeners.delete(forceUpdate);
    }, [forceUpdate]);

    return [...CspBlockedUrls] as const;
}
