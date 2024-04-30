package env

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/yyewolf/rwbyadv3/internal/values"
)

type Config struct {
	Mode values.Mode `env:"MODE" envDefault:"unset"`

	// Database
	Database struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		Port     string `env:"PORT" envDefault:"5432"`
		User     string `env:"USER" envDefault:""`
		Pass     string `env:"PASS" envDefault:""`
		Database string `env:"DATABASE" envDefault:"rwby"`

		SchemaFile       string `env:"SCHEMA_FILE" envDefault:"/sql/schema.sql"`
		MigrationsFolder string `env:"MIGRATIONS_FOLDER" envDefault:"/sql/migrations"`
	} `envPrefix:"DB_"`

	// Discord
	Discord struct {
		Token string       `env:"TOKEN" envDefault:""`
		AppID snowflake.ID `env:"APP_ID" envDefault:""`
	} `envPrefix:"DISCORD_"`

	// Github
	Github struct {
		Token      string `env:"TOKEN" envDefault:""`
		Username   string `env:"USERNAME" envDefault:""`
		Repository string `env:"REPOSITORY" envDefault:""`
	} `envPrefix:"GITHUB_"`
}

func Get() Config {
	return cfg
}
