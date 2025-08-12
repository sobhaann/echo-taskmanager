package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sobhaann/echo-taskmanager/models"
)

func CreateToken(user models.User) (string, error) {
	godotenv.Load()
	secret := []byte(os.Getenv("JWT_SECRET"))
	expMinutes, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if expMinutes == 0 {
		expMinutes = 60
	}

	claims := jwt.MapClaims{
		"user_id":      user.ID,
		"name":         user.UserName,
		"phone_number": user.PhoneNumber,
		"exp":          time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)

}
