package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ujjwal405/google_login/helper"
	model "github.com/ujjwal405/google_login/models"
)

var (
	authType      = "bearer"
	errInvalid    = errors.New("invalid format")
	errNotSupport = errors.New("doesn't support this type of authorization")
	errCookie     = errors.New("cannot get cookie")
)

var claim = model.ContextKey("claim")

func RecoveryHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					res := e.Error()
					log.Println(res)
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					w.Write([]byte(res))
					return
				}
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func AuthHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.URL.Path {
		case "/main", "/edit", "/save":
			cookie, err := r.Cookie("Authorization")
			if err == http.ErrNoCookie {
				panic(errCookie)
			}
			cookievalue := cookie.Value

			field := strings.Fields(cookievalue)
			if len(field) < 2 {
				panic(errInvalid)
			}
			authtype := strings.ToLower(field[0])

			if authtype != authType {
				panic(errNotSupport)
			}
			token := field[1]
			claims, err := helper.ValidateToken(token)
			if err != nil {
				panic(err)
			}
			Claims := *claims
			ctx := context.WithValue(r.Context(), claim, Claims)
			next.ServeHTTP(w, r.WithContext(ctx))

		default:
			next.ServeHTTP(w, r)
		}
	}
}
