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

		App struct {
			ClientID     string `env:"CLIENT_ID" envDefault:""`
			ClientSecret string `env:"CLIENT_SECRET" envDefault:""`
			BaseURI      string `env:"BASE_URI" envDefault:""`
		} `envPrefix:"APP_"`
	} `envPrefix:"DISCORD_"`

	// Github
	Github struct {
		Token      string `env:"TOKEN" envDefault:""`
		Username   string `env:"USERNAME" envDefault:""`
		Repository string `env:"REPOSITORY" envDefault:""`

		App struct {
			ClientID     string `env:"CLIENT_ID" envDefault:""`
			ClientSecret string `env:"CLIENT_SECRET" envDefault:""`
			BaseURI      string `env:"BASE_URI" envDefault:""`
		} `envPrefix:"APP_"`
	} `envPrefix:"GITHUB_"`

	// Rbmq
	Rbmq struct {
		Host string `env:"HOST" envDefault:"localhost"`
		Port string `env:"PORT" envDefault:"5672"`
		User string `env:"USER" envDefault:""`
		Pass string `env:"PASS" envDefault:""`

		Jobs struct {
			DLExchange string `env:"DLEXCHANGE" envDefault:"jobs"`
			Exchange   string `env:"EXCHANGE" envDefault:"jobs"`
			Queue      string `env:"QUEUE" envDefault:"jobs"`
		} `envPrefix:"JOBS_"`
	} `envPrefix:"RBMQ_"`

	// Temporal
	Temporal struct {
		Host      string `env:"HOST" envDefault:"localhost"`
		Port      string `env:"PORT" envDefault:"7233"`
		TaskQueue string `env:"TASK_QUEUE" envDefault:"worker-dev"`
	} `envPrefix:"TEMPORAL_"`

	// Web
	Web struct {
		Port string `env:"PORT" envDefault:"8080"`
	} `envPrefix:"WEB_"`

	// App parameters
	App struct {
		CardsLocation string `env:"CARDS_LOCATION" envDefault:"/cards/yml"`
		BotColor      int    `env:"BOT_COLOR" envDefault:"3859607"`
		BackpackSize  int    `env:"BACKPACK_SIZE" envDefault:"20"`
		BaseURI       string `env:"BASE_URI" envDefault:""`
	} `envPrefix:"APP_"`
}

func Get() *Config {
	return &cfg
}
