package enums

type AuthStateTypes string

const (
	DiscordLogin    AuthStateTypes = "discord_login"
	GithubCheckStar AuthStateTypes = "discord_check_star"
)

// Values provides list valid values for Enum.
func (AuthStateTypes) Values() (kinds []string) {
	for _, s := range []AuthStateTypes{DiscordLogin, GithubCheckStar} {
		kinds = append(kinds, string(s))
	}
	return
}
