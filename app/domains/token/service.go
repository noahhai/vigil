package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

type Service interface {
	Generate(username, email string) (tkn string, err error)
	GenerateAgentToken() (tkn string)
}

type tokenService struct {
	signingKey []byte
}

func NewTokenService() Service {
	return &tokenService{
		signingKey: []byte(os.Getenv("SIGNING_KEY")),
	}
}

func (t *tokenService) Generate(username, email string) (tkn string, err error) {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["admin"] = false
	claims["username"] = username
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tkn, err = token.SignedString(t.signingKey)
	return
}

func (t *tokenService) GenerateAgentToken() (tkn string) {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
