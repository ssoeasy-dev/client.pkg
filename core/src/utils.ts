// Проверка на браузерное окружение
const isBrowser =
  typeof window !== "undefined" && typeof window.crypto !== "undefined";

function base64URLEncode(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer);
  let binary = "";
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary)
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
}

export async function generateCodeVerifier(): Promise<string> {
  if (!isBrowser) {
    throw new Error(
      "generateCodeVerifier can only be used in browser environment",
    );
  }
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return base64URLEncode(array.buffer);
}

export async function generateCodeChallenge(verifier: string): Promise<string> {
  if (!isBrowser) {
    throw new Error(
      "generateCodeChallenge can only be used in browser environment",
    );
  }
  const encoder = new TextEncoder();
  const data = encoder.encode(verifier);
  const digest = await crypto.subtle.digest("SHA-256", data);
  return base64URLEncode(digest);
}

export function generateState(): string {
  if (!isBrowser) {
    throw new Error("generateState can only be used in browser environment");
  }
  // Используем crypto.randomUUID, если доступен, иначе fallback
  if (crypto.randomUUID) {
    return crypto.randomUUID();
  }
  // Простой генератор UUID v4 (для старых браузеров)
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}
