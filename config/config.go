package config

type Configuration struct {
	AppName   string `env:"APP_NAME" env-default:"GO_ECHO_GORM_ASSO"`
	AppEnv    string `env:"APP_ENV" env-default:"DEV"`
	Port      string `env:"MY_APP_PORT" env-default:"8080"`
	Host      string `env:"HOST" env-default:"localhost"`
	LogLevel  string `env:"LOG_LEVEL" env-default:"ERROR"`
	JwtSecret string `env:"JWT_SECRET" env-default:"debalina"`
}
