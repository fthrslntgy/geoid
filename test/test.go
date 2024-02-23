package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/fthrslntgy/geoid"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/oauth2"
)

var config *oauth2.Config
var user_info_endpoint string

var provider_files = []string{"google_test_provider", "yahoo_test_provider"}

func main() {
	var err error
	config, user_info_endpoint, err = geoid.OAuthConfig(provider_files[0], false)
	if err != nil {
		panic(err)
	}
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/login", Login)
	app.Get("/callback", Callback)
	app.ListenTLS(":8080", "cert.pem", "key.pem")
}

func Login(c *fiber.Ctx) error {
	url := config.AuthCodeURL("randomstate")
	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return c.JSON(url)
}

func Callback(c *fiber.Ctx) error {
	state := c.Query("state")
	if state != "randomstate" {
		return fmt.Errorf("states do not match")
	}
	code := c.Query("code")
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("code-token exchange failed with error: %s", err)
	}
	req, err := http.NewRequest("GET", user_info_endpoint, nil)
	if err != nil {
		return fmt.Errorf("user data fetch failed with error: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("user data fetch failed with error: %s", err)
	}
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("json parsing failed with error: %s", err)
	}

	return c.SendString(string(userData))
}
