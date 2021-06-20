package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/kangana1024/goadmin-workshop/database"
	"github.com/kangana1024/goadmin-workshop/models"
	"golang.org/x/crypto/bcrypt"
)

func Home(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello Fiber ✨",
	})
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

	if err := c.BodyParser(&req); err != nil {
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

type SignINRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignINResponse struct {
	ID  int64
	JWT string `json:"jwt"`
}

func Signin(c *fiber.Ctx) error {
	var req SignINRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email=?", req.Email).First(&user)
	if user.ID == 0 {
		c.Status(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User Not Found!",
		})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		c.Status(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Permission Denied!",
		})
	}

	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Permission Denied!",
			"error":   fmt.Sprintf("%+v : %+v", err, jwtToken),
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(&SignINResponse{
		ID:  int64(user.ID),
		JWT: token,
	})
}

type Claims struct {
	jwt.StandardClaims
}

func User(c *fiber.Ctx) error {

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Permission Denied!",
		})
	}
	claim := token.Claims.(*Claims)
	id := claim.Issuer
	var user models.User
	database.DB.Where("id=?", id).First(&user)

	return c.JSON(user)
}

func SignOut(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}
