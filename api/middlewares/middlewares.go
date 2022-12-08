package middlewares

import (
	"encoding/json"
	"errors"
	"github.com/gmarshall142/services/api/auth"
	"github.com/gmarshall142/services/api/responses"
	"log"
	"net/http"
)

const ctxTokenKey = "Auth0Token"

type message struct {
	Message string `json:"message"`
}

func sendMessage(rw http.ResponseWriter, data *message) {
	rw.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Print("json conversion error", err)
		return
	}
	_, err = rw.Write(bytes)
	if err != nil {
		log.Print("http response write error", err)
	}
}

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}

func SetMiddlewareValidToken(next http.HandlerFunc) http.Handler {
	return EnsureValidToken(next)
}

func ValidateTokenAndScope(next http.HandlerFunc, scope string) http.Handler {
	return EnsureValidTokenAndScope(next, scope)
}

// validateToken middleware verifies a valid Auth0 JWT token being present in the request.
//func ValidateToken(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
//		token, err := extractToken(req)
//		if err != nil {
//			fmt.Printf("failed to parse payload: %s\n", err)
//			rw.WriteHeader(http.StatusUnauthorized)
//			sendMessage(rw, &message{err.Error()})
//			return
//		}
//		ctxWithToken := context.WithValue(req.Context(), ctxTokenKey, token)
//		next.ServeHTTP(rw, req.WithContext(ctxWithToken))
//	})
//}
