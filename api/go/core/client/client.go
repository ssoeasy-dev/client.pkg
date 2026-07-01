package client

import (
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// Client — HTTP-клиент к auth.api.
// Создайте один раз и переиспользуйте во всём приложении.
type Client struct {
	baseURL        url.URL
	httpClient     *http.Client
	serviceID      uuid.UUID
}

// NewClient создаёт Client.
//   - env — окружение, в котором запускается приложение.
//     Валидные значения: EnvDevelopment, EnvProduction, EnvLocal.
//     С EnvLocal нужно обязательно передать опцию WithBaseURL.
//   - opts — опции сборки. Включает значения: WithHttpTimeout, WithBaseURL.
func NewClient(cfg Config) (*Client, error) {
	url, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	c := &Client{
		baseURL: *url,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		serviceID:      cfg.ServiceID,
	}

	return c, nil
}
