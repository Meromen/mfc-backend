package preferences

import (
	"github.com/kelseyhightower/envconfig"
)

type Preferences struct {
	LogLevel                int           `envconfig:"LOG_LEVEL" required:"true"`
	LogAsJSON               bool          `envconfig:"LOG_AS_JSON" required:"true"`
	PostgresUrl             string        `envconfig:"POSTGRES_URL" required:"true"`
	ServerAddress           string        `envconfig:"SERVER_ADDRESS" required:"false"`
}

func Get() (*Preferences, error) {
	var p Preferences
	err := envconfig.Process("", &p)
	return &p, err
}
