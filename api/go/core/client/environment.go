package client

import ssoeasy "github.com/ssoeasy-dev/client.pkg/api/go/core"

// Окружение сборки
type Environment string

const (
	EnvDevelopment Environment = "Development" // Окружение разработки. Не для продакшена
	EnvProduction  Environment = "Production"  // Окружение продакшена. Только для рабочей версии
	EnvEnterprise  Environment = "Enterprise"  // Окружение собственного инстанса. Только для enterprise клиентов
)

func (e Environment) baseUrl() (string, error) {
	switch e {
	case EnvDevelopment:
		return baseUrlDevelopment, nil
	case EnvProduction:
		return baseUrlProduction, nil
	case EnvEnterprise:
		return "", nil
	}
	return "", ssoeasy.NewError("unknown enviroment")
}
