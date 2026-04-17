const STORAGE_PREFIX = "sso_auth_";

export interface TempAuthData {
  verifier: string;
  redirectTo: string;
}

function getStorage(): Storage | null {
  if (typeof localStorage !== "undefined") {
    return localStorage;
  }
  return null;
}

export function saveTempData(state: string, data: TempAuthData): void {
  const storage = getStorage();
  if (!storage) {
    console.warn("localStorage not available");
    return;
  }
  storage.setItem(STORAGE_PREFIX + state, JSON.stringify(data));
}

export function getTempData(state: string): TempAuthData | null {
  const storage = getStorage();
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
  const storage = getStorage();
  if (!storage) return;
  storage.removeItem(STORAGE_PREFIX + state);
}
