package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/raihan2bd/vidverse/config"
	"github.com/raihan2bd/vidverse/models"
)

// Decode the token
func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// Validate the token
func ValidateToken(claims jwt.MapClaims) bool {
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return false
	}

	// check the user_id from the token as well
	if claims["user_id"] != 0 {
		return false
	}

	return true
}

// Validate user_id and convert to uint
func ValidateAndGetUserByID(app *config.Application, id any) (*models.User, error) {
	userID, ok := id.(uint)
	if !ok {
		return nil, errors.New("invalid User")
	}

	user, err := app.DBMethods.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("invalid User")
	}

	if user.ID <= 0 {
		return nil, errors.New("invalid User")
	}
	user.Password = ""

	return user, nil
}
