package echo

import "github.com/ssoeasy-dev/client.pkg/api/go/core/v2/client"

type Config struct {
	Cookie CookieConfig
	Client client.Config
}

type CookieConfig struct {
	Domain string
	MaxAge int
	Path   string
}
