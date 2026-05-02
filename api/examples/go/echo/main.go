package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ssoeasy-dev/client.pkg/api/go/core/client"
	"github.com/ssoeasy-dev/client.pkg/api/go/echo"
)

func main() {
	// Инициализация клиента SSO. От env зависит среда ssoeasy.
	// EnvProduction используется для продакшен среды ssoeasy.
	// EnvDevelopment используется для среды разработки и ограничивает количество запросов.
	// EnvEnterprise используется для использования собственного инстанса ssoeasy. Требует опции ssoeasyClient.WithBaseURL()
	client, err := client.NewClient(client.EnvDevelopment)
	if err != nil {
		log.Fatal(err)
	}

	// ID права, которое требуется для доступа. Выпускается в админке ssoeasy
	permID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	e := echo.New()

	// --- Использование ---

	// 1. Публичный эндпоинт — middleware ExtractUser добавляет пользователя в контекст, если токен передан
	e.GET("/public", handler, ssoeasy.ExtractUser())

	// 2. Защищённый эндпоинт — middleware CheckPermission проверяет доступ пользователя в системе ssoeasy по токену и добавляет пользователя в контекст
	e.GET("/protected", handler, ssoeasy.CheckPermission(client, permID))

	e.Logger.Fatal(e.Start(":8080"))
}

func handler(c echo.Context) error {
	// ssoeasyMiddleware.GetUser получает пользователя из контекста.
	user := ssoeasy.GetUser(c)
	// Проверка нужна, если используется ssoeasyMiddleware.ExtractUser напрямую.
	if user == nil {
		return c.String(http.StatusOK, "Anonymous")
	}
	return c.String(http.StatusOK, "Hello, "+user.UserID.String())
}
