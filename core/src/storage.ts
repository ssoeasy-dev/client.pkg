const STORAGE_PREFIX = "sso_auth_";

export interface TempAuthData {
  verifier: string;
  redirectTo: string;
}

function getSessionStorage(): Storage | null {
  if (typeof sessionStorage !== "undefined") {
    return sessionStorage;
  }
  return null;
}

export function saveTempData(state: string, data: TempAuthData): void {
  const storage = getSessionStorage();
  if (!storage) {
    console.warn(
      "sessionStorage not available, authentication flow may not work",
    );
    return;
  }
  storage.setItem(STORAGE_PREFIX + state, JSON.stringify(data));
}

export function getTempData(state: string): TempAuthData | null {
  const storage = getSessionStorage();
  if (!storage) return null;
  const raw = storage.getItem(STORAGE_PREFIX + state);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as TempAuthData;
  } catch {
    return null;
  }
}

export function removeTempData(state: string): void {
  const storage = getSessionStorage();
  if (!storage) return;
  storage.removeItem(STORAGE_PREFIX + state);
}
