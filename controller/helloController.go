package controller

import (
	"strconv"
	"time"

	"example.com/go/auth/database"
	"example.com/go/auth/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = "secret"

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.ErrBadRequest)

	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	database.DBCon.Create(&user)
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return msg(c, fiber.StatusBadRequest, "Bad Request")

	}

	var user models.User

	database.DBCon.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		return msg(c, fiber.StatusNotFound, "User not found")

	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return msg(c, fiber.StatusBadRequest, "Incorrect password")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		return msg(c, fiber.StatusInternalServerError, "could not login")

	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return msg(c, fiber.StatusAccepted, "login success")
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return msg(c, fiber.StatusUnauthorized, "Unauathenticated")
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DBCon.Where("id=?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return msg(c, fiber.StatusAccepted, "logout success")

}

// private function

func msg(c *fiber.Ctx, code int, message string) error {
	c.Status(code)
	return c.JSON(fiber.Map{
		"message": message,
	})

}
