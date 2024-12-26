package config

type Config struct {
	SentryDSN      string `env:"SENTRY_DSN"`
	SendgridApiKey string `env:"SENDGRID_API_KEY"`
	S3SecretKey    string `env:"S3_SECRET_KEY"`
	S3AccessKey    string `env:"S3_ACCESS_KEY"`
	S3Endpoint     string `env:"S3_ENDPOINT"`
	Env            string `env:"ENV"`
	SenderEmail    string `env:"SENDER_EMAIL"`
	BaseURL        string `env:"BASE_URL"`
}

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)
