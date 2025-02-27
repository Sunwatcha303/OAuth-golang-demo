package servers

import (
	"log"

	"github.com/Sunwatcha303/OAuth-golang-demo/configs"
	"github.com/Sunwatcha303/OAuth-golang-demo/pkg/databases"
	"github.com/Sunwatcha303/OAuth-golang-demo/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	App *fiber.App
	Cfg *configs.Configs
	Db  *databases.Database
}

func NewServer(cfg *configs.Configs, db *sqlx.DB, redis *redis.Client) *Server {
	return &Server{
		App: fiber.New(),
		Cfg: cfg,
		Db: &databases.Database{
			PostgreSQL: db,
			Redis:      redis,
		},
	}
}

func (s *Server) Start() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	if err := s.MapHandlers(); err != nil {
		log.Fatalln(err.Error())
		panic(err.Error())
	}

	fiberConnURL, err := utils.ConnectionUrlBuilder("fiber", s.Cfg)
	if err != nil {
		log.Fatalln(err.Error())
		panic(err.Error())
	}

	host := s.Cfg.App.Host
	port := s.Cfg.App.Port
	log.Printf("server has been started on %s:%s ⚡", host, port)

	if err := s.App.Listen(fiberConnURL); err != nil {
		log.Fatalln(err.Error())
		panic(err.Error())
	}
}
