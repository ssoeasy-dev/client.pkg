export const BASE_URL = {
    development: "https://auth.dev.ssoeasy.ru",
    production: "https://auth.ssoeasy.ru",
}

export const AUTH_PAGE = {
    development: "https://auth.dev.ssoeasy.ru/login",
    production: "https://auth.ssoeasy.ru/login",
}

export const AUTHORIZE_PATH = {
    development: "/api/v1/auth/authorize/:serviceId",
    production: "/api/v1/auth/authorize/:serviceId",
}

export const LOGOUT_PATH = {
    development: "/api/v1/auth/logout",
    production: "/api/v1/auth/logout",
}

export const ME_PATH = {
    development: "/api/v1/auth/me",
    production: "/api/v1/auth/me",
}
