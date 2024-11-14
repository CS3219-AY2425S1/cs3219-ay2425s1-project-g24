"use client";

const KEY: string = "TOKEN";

export function setToken(token: string): void {
    document.cookie = `${KEY}=${token}`
}

export function getToken(): string {
    const keyValue = document.cookie.split("; ").find(kv => kv.startsWith(`${KEY}=`))
    if (keyValue == undefined) {
        return "";
    }

    const [_, value] = keyValue.split("=");

    if (value == undefined) {
        return "";
    }

    return value;
}

export function deleteToken() {
    document.cookie = `${KEY}=nothinghere;expires=0`;
}