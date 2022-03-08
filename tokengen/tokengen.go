package tokengen

import (
	"os"

	"apricate/log"

	"github.com/golang-jwt/jwt"
)

// Generates a new token based on username and apricate_access_secret
func GenerateToken(username string) (string, error) {
	// Creating access token
	// Set claims for jwt
	atClaims := jwt.MapClaims{}
	atClaims["username"]=username
	// Use signing method HS256
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	log.Debug.Printf("Got APRICATE_ACCESS_SECRET:\n%s", os.Getenv("APRICATE_ACCESS_SECRET"))
	// Generate token using apricate_access_secret
	token, err := at.SignedString([]byte(os.Getenv("APRICATE_ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
