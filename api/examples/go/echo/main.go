package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	ssoeasyClient "github.com/ssoeasy-dev/client.pkg/api/go/core/v2/client"
	ssoeasyDto "github.com/ssoeasy-dev/client.pkg/api/go/core/v2/dto"
	ssoeasyEcho "github.com/ssoeasy-dev/client.pkg/api/go/echo/v2"
)

func main() {
	permissionID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	serviceID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	e := echo.New()

	authHandler, err := ssoeasyEcho.NewHandler(ssoeasyEcho.Config{
		Cookie: ssoeasyEcho.CookieConfig{
			Domain: "example.com",
			MaxAge: 8 * 60 * 60 * 1000,
		},
		Client: ssoeasyClient.Config{
			BaseURL: ssoeasyClient.ProductionBaseUrl,
			ServiceID: serviceID,
		},
	})
	if err != nil {
		panic(err)
	}

	cb := func (ctx context.Context, p ssoeasyDto.Payload) error { return nil }

	authGroup := e.Group("/auth")
	authGroup.POST("/authorize", authHandler.Authorize(cb))
	authGroup.DELETE("/logout", authHandler.Logout(cb))
	authGroup.GET("/me", authHandler.Me(cb))
	authGroup.PATCH("/refresh", authHandler.Refresh(cb))

	
	e.GET("/public", handler)
	e.GET("/protected", handler, authHandler.Check(permissionID, cb))

	e.Logger.Fatal(e.Start(":8080"))
}

func handler(c echo.Context) error {
	user, err := ssoeasyDto.PayloadFromContext(c.Request().Context())
    if err != nil {
        return err
    }
	return c.String(http.StatusOK, "Hello, "+user.UserID.String())
}
