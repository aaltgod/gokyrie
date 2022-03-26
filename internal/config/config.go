package config

type Config struct {
	Teams []Team
}

type Team struct {
	Name, IP string
}

func NewConfig(teams ...Team) *Config {
	cfg := new(Config)

	cfg.Teams = teams

	return cfg
}
