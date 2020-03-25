package main

import (
	"fiber-jwt/models"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"github.com/raymayemir/jwt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Fiber-JWT example v1.0.0")
	e := godotenv.Load()

	if e != nil {
		fmt.Println(e)
	}
	port := os.Getenv("port")

	app := fiber.New()

	cfg := jwt.Config{
		NotAuth:       []string{"/", "/login"},
		TokenPassword: os.Getenv("jwt-secret"),
	}

	app.Use(jwt.New(cfg))

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello Guest!")
	})

	app.Post("/login", func(c *fiber.Ctx) {
		user := &models.User{}
		if err := c.BodyParser(&user); err != nil {
			// Handle error
			fmt.Println(err.Error())
			_ = c.JSON("error")
			return
		}

		token := sign(user.ID, os.Getenv("jwt-secret"))
		_ = c.JSON(token)

	})

	app.Get("/user", func(c *fiber.Ctx) {

		id := getUserId(c.Get(fiber.HeaderAuthorization), os.Getenv("jwt-secret"))
		if id == "empty" {
			_ = c.JSON("error id not found")
			return
		}

		_ = c.JSON(fmt.Sprintf(`Hello user! Your ID is: %s`, id))
	})

	_ = app.Listen(port)
}

func sign(id uint, secret string) string {
	tk := &models.Token{UserID: fmt.Sprint(id)}
	token := jwtgo.NewWithClaims(jwtgo.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(secret))
	fmt.Println(tokenString)
	return tokenString
}
func getUserId(tokenHeader, secret string) string {
	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) == 2 {
		tokenPart := splitted[1] // Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwtgo.ParseWithClaims(tokenPart, tk, func(token *jwtgo.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err == nil && token.Valid { // Malformed token, returns with http code 403 as usual

			return tk.UserID
		}
	}

	return "empty"
}
