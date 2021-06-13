package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kangana1024/goadmin-workshop/database"
	"github.com/kangana1024/goadmin-workshop/models"
	"golang.org/x/crypto/bcrypt"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("OK ✨")
}

type RegisterRequest struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}

	if req.Password != req.PasswordConfirm {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Password do not match!",
		})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return err
	}

	var user models.User
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.Password = password

	database.DB.Create(&user)

	return c.JSON(user)
}
