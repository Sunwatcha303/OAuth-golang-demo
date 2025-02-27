package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Sunwatcha303/OAuth-golang-demo/configs"
	"github.com/Sunwatcha303/OAuth-golang-demo/modules/servers"
	"github.com/Sunwatcha303/OAuth-golang-demo/pkg/databases"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		panic(err.Error())
	}
	cfg := new(configs.Configs)
	cfg.App.Host = os.Getenv("FIBER_HOST")
	cfg.App.Port = os.Getenv("FIBER_PORT")

	cfg.PostgreSQL.Host = os.Getenv("DB_HOST")
	cfg.PostgreSQL.Port = os.Getenv("DB_PORT")
	cfg.PostgreSQL.Protocol = os.Getenv("DB_PROTOCOL")
	cfg.PostgreSQL.Username = os.Getenv("DB_USERNAME")
	cfg.PostgreSQL.Password = os.Getenv("DB_PASSWORD")
	cfg.PostgreSQL.Database = os.Getenv("DB_DATABASE")

	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	cfg.Redis.Database, _ = strconv.Atoi(os.Getenv("REDIS_DATABASE"))

	cfg.OAuth.ClientID = os.Getenv("CLIENT_ID")
	cfg.OAuth.ClientSecret = os.Getenv("CLIENT_SECRET")
	cfg.OAuth.RedirectUri = os.Getenv("REDIRECT_URI")

	cfg.Jwt.SecretKey = os.Getenv("JWT_SECRET_KEY")

	db, err := databases.NewPostgreSQLDBConnection(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	redis, err := databases.NewRedisConnection(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer db.Close()
	defer redis.Close()

	s := servers.NewServer(cfg, db, redis)
	s.Start()
}
