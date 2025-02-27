package servers

import (
	_usersHttp "github.com/Sunwatcha303/OAuth-golang-demo/modules/users/controllers"
	_usersRepository "github.com/Sunwatcha303/OAuth-golang-demo/modules/users/repositories"
	_usersUsecase "github.com/Sunwatcha303/OAuth-golang-demo/modules/users/usecases"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) MapHandlers() error {

	v1 := s.App.Group("/v1")
	usersGroup := v1.Group("/users")
	usersRepository := _usersRepository.NewUsersRepository(s.Db)
	usersUsecase := _usersUsecase.NewUsersUsecase(usersRepository, s.Cfg)
	_usersHttp.NewUsersController(usersGroup, usersUsecase)

	// End point not found response
	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     "error, end point not found",
		})
	})

	return nil
}
