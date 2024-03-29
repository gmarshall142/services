package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var publicKey *rsa.PublicKey

func init() {
	key, _ := os.ReadFile("certs/dev-rn77drwl.pem")
	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	if err != nil {
		fmt.Errorf("error verifying token: %v", err)
	}
	publicKey = parsedKey
}

func CreateToken(user_id uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

func TokenValid(r *http.Request) error {
	_, err := getToken(r)
	return err
}

func TokenAndPermValid(r *http.Request, perm string) error {
	token, err := getToken(r)
	if err != nil {
		return err
	}
	var fnd = false
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// TODO: find a better way to get permissions from {}interface
		for key, v := range claims {
			if key == "permissions" {
				fmt.Printf("v: %v\n", v)
				vals := fmt.Sprintf("%v", v)
				vals2 := strings.Split(strings.Trim(strings.Trim(vals, "["), "]"), " ")
				fmt.Println(vals2)
				for _, str := range vals2 {
					if str == perm {
						fnd = true
					}
				}
			}
		}
	}
	if !fnd {
		return fmt.Errorf("permission not found")
	}
	return nil
}

func getToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// Pretty display the claims nicely in the terminal
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
