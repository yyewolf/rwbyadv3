package env

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/yyewolf/rwbyadv3/internal/values"
)

type Config struct {
	Mode values.Mode `env:"MODE" envDefault:"unset"`

	// Database
	Mongo struct {
		Host       string `env:"HOST" envDefault:"localhost"`
		Port       string `env:"PORT" envDefault:"27017"`
		User       string `env:"USER" envDefault:""`
		Pass       string `env:"PASS" envDefault:""`
		Database   string `env:"DATABASE" envDefault:"rcbs"`
		Additional string `env:"ADDITIONAL" envDefault:""`
	} `envPrefix:"MONGO_"`

	// Discord
	Discord struct {
		Token          string        `env:"TOKEN" envDefault:""`
		AppID          string        `env:"APP_ID" envDefault:""`
		AppIDSnowflake discord.AppID `env:"-"`
	} `envPrefix:"DISCORD_"`
}

func Get() Config {
	return cfg
}
