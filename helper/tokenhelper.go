package helper

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Signedetails struct {
	User_id string
	jwt.StandardClaims
}

func GenerateToken(uid string, duration time.Duration) (token string, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return "", err
	}
	secret := os.Getenv("SECRET_KEY")
	claims := Signedetails{
		User_id: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(duration)).Unix(),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateToken(SignedToken string) (*Signedetails, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	var SECRET_KEY string = os.Getenv("SECRET_KEY")

	token, err := jwt.ParseWithClaims(
		SignedToken,
		&Signedetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Signedetails)
	if !ok {

		err = errors.New("token invalid")
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}
	return claims, err

}
func HashPassword(userpassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(userpassword), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	check := true
	msg := ""
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		msg = "password incorect"
		check = false
	}
	return check, msg
}

func GenerateRandom() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uid := newUUID.String()
	return uid, nil
}
func SetCookie(w http.ResponseWriter, uid string) {
	cookie := &http.Cookie{
		Name:     "uid",
		Value:    uid,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

}
func DeleteCookie(w http.ResponseWriter, uid string) {
	cookie := &http.Cookie{
		Name:   "uid",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

}
func SetAuthCookie(w http.ResponseWriter, auth string) {
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    auth,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

}
func DeleteAuthCookie(w http.ResponseWriter, auth string) {
	cookie := &http.Cookie{
		Name:   "Authorization",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

}
