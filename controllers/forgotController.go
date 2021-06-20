package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"

	"github.com/gofiber/fiber/v2"
	"github.com/kangana1024/goadmin-workshop/database"
	"github.com/kangana1024/goadmin-workshop/models"
	"golang.org/x/crypto/bcrypt"
)

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

func Forgot(c *fiber.Ctx) error {
	var req ForgotPasswordReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	token := RandStringRunes(12)

	passwordReset := models.PasswordReset{
		Email: req.Email,
		Token: token,
	}

	database.DB.Create(&passwordReset)

	from := "admin@examplo.com"

	url := "/reset/" + token

	message := []byte(fmt.Sprintf("Click : <a href='%s'>reset password</a>", url))

	err := smtp.SendMail("0.0.0.0:1025", nil, from, []string{
		req.Email,
	}, []byte(message))

	if err != nil {
		return nil
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

type ResetPasswordReq struct {
	Password        string
	PasswordConfirm string `json:"password_confirm"`
	Token           string
}

func Reset(c *fiber.Ctx) error {
	var req ResetPasswordReq

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if req.Password != req.PasswordConfirm {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Password not match!",
		})
	}

	var passwordReset models.PasswordReset
	if err := database.DB.Where("token=?", req).Last(&passwordReset); err.Error != nil {
		c.Status(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid Token!",
		})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	if err != nil {
		return err
	}

	database.DB.Model(&models.User{}).Where("email=?", passwordReset.Email).Update("password", password)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
