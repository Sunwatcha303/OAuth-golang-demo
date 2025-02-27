package databases

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	PostgreSQL *sqlx.DB
	Redis      *redis.Client
}
