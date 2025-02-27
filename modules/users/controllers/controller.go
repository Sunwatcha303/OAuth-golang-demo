package controllers

import (
	"strings"

	"github.com/Sunwatcha303/OAuth-golang-demo/modules/entities"
	"github.com/Sunwatcha303/OAuth-golang-demo/modules/users/usecases"
	"github.com/gofiber/fiber/v2"
)

type UsersController struct {
	userUsecase *usecases.UserUsecase
}

func NewUsersController(r fiber.Router, userUsecase *usecases.UserUsecase) {
	controller := &UsersController{
		userUsecase: userUsecase,
	}
	r.Get("/", controller.Health)
	r.Get("/login", controller.LoginGoogle)
	r.Get("/authentication", controller.Authentication)
	r.Post("/oauth", controller.OAuth)

}

func (h *UsersController) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Server is running",
	})
}

func (h *UsersController) Authentication(c *fiber.Ctx) error {
	authHeader := strings.Split(c.Get("Cookie"), "token=")[1]
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Unauthorized",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Authorization token is required",
			"result":      nil,
		})
	}
	token, err := h.userUsecase.VerifyAndExtractToken(authHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Unauthorized",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Invalid token or failed to decode",
			"error":       err.Error(),
		})
	}

	newToken, err := h.userUsecase.GetToken(token)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "StatusInternalServerError",
			"status_code": fiber.StatusInternalServerError,
			"message":     "Failed to decode response",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"data":        newToken,
	})
}

func (h *UsersController) LoginGoogle(c *fiber.Ctx) error {
	url := h.userUsecase.GetUrlOAuth()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"data":        url,
	})
}

func (h *UsersController) OAuth(c *fiber.Ctx) error {
	var authCode entities.OAuthRequest
	if err := c.BodyParser(&authCode); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}
	token, err := h.userUsecase.GetNewToken(authCode.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "StatusInternalServerError",
			"status_code": fiber.StatusInternalServerError,
			"message":     "Creation token not success",
			"error":       err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Login success",
		"data":        token,
	})
}
