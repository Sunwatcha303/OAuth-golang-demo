package configs

type Configs struct {
	PostgreSQL PostgreSQL
	Redis      Redis
	App        Fiber
	OAuth      OAuth
	Jwt        JWT
}

type Fiber struct {
	Host string
	Port string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Protocol string
	Username string
	Password string
	Database string
	SSLMode  string
}

type OAuth struct {
	ClientID     string
	ClientSecret string
	RedirectUri  string
}

type Redis struct {
	Host     string
	Port     string
	Database int
}

type JWT struct {
	SecretKey string
}
