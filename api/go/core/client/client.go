package client

import (
	"net/http"
	"strings"
	"time"
)

// Client — HTTP-клиент к auth.api.
// Создайте один раз и переиспользуйте во всём приложении.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// ClientOption позволяет настроить Client.
type clientOption func(*Client) error

// NewClient создаёт Client.
//   - env — окружение, в котором запускается приложение.
//     Валидные значения: EnvDevelopment, EnvProduction, EnvLocal.
//     С EnvLocal нужно обязательно передать опцию WithBaseURL.
//   - opts — опции сборки. Включает значения: WithHttpTimeout, WithBaseURL.
func NewClient(env Environment, opts ...clientOption) (*Client, error) {
	baseURL, err := env.baseUrl()
	if err != nil {
		return nil, err
	}
	c := &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithHttpTimeout опция для изменения таймаута Client.
//   - timeout — таймаут в секундах.
func WithHttpTimeout(timeout time.Duration) clientOption {
	return func(c *Client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}

// WithBaseURL опция для изменения базового url Client.
//   - baseURL — url сервера аутентификации.
func WithBaseURL(baseURL string) clientOption {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}
