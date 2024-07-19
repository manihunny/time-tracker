package config

type Config struct {
	ServerHost string `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort string `env:"SERVER_PORT" envDefault:"80"`

	DBDriver        string `env:"DB_DRIVER"`
	DBHost          string `env:"DB_HOST"`
	DBPort          string `env:"DB_PORT"`
	DBName          string `env:"DB_NAME"`
	DBUser          string `env:"DB_USER"`
	DBPassword      string `env:"DB_PASSWORD"`
	DBMigrationsDir string `env:"DB_MIGRATIONS_DIR"`
}
