package tokenconfig

import (
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func LoginConfig() (conf *oauth2.Config, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return
	}
	client_id := os.Getenv("client_id")
	client_secret := os.Getenv("client_secret")

	conf = &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		RedirectURL:  "http://localhost:9090/logincallback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return

}
