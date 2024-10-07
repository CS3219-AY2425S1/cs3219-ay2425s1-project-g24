const KEY: string = "TOKEN";

export function setToken(token: string): void {
    window.localStorage[KEY] = token
}

export function getToken(): string {
    return KEY in window.localStorage ? window.localStorage[KEY] : ""
}

export function deleteToken() {
    if (KEY in window.localStorage) {
        window.localStorage.removeItem(KEY);
    }
}