package infrastructure

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/orkungursel/hey-taxi-location-api/config"
	"github.com/orkungursel/hey-taxi-location-api/internal/app"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
	"github.com/pkg/errors"
)

type TokenService struct {
	config               *config.Config
	logger               logger.ILogger
	accessTokenPublicKey *rsa.PublicKey
}

func NewTokenService(config *config.Config, logger logger.ILogger) (s *TokenService) {
	s = &TokenService{
		config: config,
		logger: logger,
	}

	s.init()

	return
}

func (s *TokenService) init() *TokenService {
	// get public key from file
	b, err := ioutil.ReadFile(s.config.Jwt.AccessTokenPublicKeyFile)
	if err != nil {
		panic(err)
	}

	atpuk, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		panic(err)
	}

	s.accessTokenPublicKey = atpuk

	return s
}

// ExtractToken extracts token from http request
func (t *TokenService) ValidateAccessTokenFromRequest(ctx context.Context, r *http.Request) (app.Claims, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("token is empty")
	}

	// remove bearer prefix
	token = token[7:]

	claims, err := t.ParseToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// ParseToken parses a token
func (t *TokenService) ParseToken(ctx context.Context, token string) (app.Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, t.provideAccessTokenPublicKey)

	if err != nil {
		return nil, err
	}

	if tkn.Method.Alg() != jwt.SigningMethodRS256.Alg() {
		return nil, errors.New("invalid algorithm")
	}

	if !tkn.Valid {
		return nil, errors.New("token is invalid")
	}

	if !claims.VerifyIssuer(t.config.Jwt.Issuer, true) {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}

// provideAccessTokenPublicKey provides access token public key to veriy token
func (t *TokenService) provideAccessTokenPublicKey(_ *jwt.Token) (interface{}, error) {
	return t.accessTokenPublicKey, nil
}
