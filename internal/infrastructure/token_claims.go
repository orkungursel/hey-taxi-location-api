package infrastructure

import "github.com/golang-jwt/jwt"

type Claims struct {
	Role string `json:"role,omitempty"`
	jwt.StandardClaims
}

func (c *Claims) GetSubject() string {
	return c.Subject
}

func (c *Claims) GetRole() string {
	return c.Role
}

func (c *Claims) GetIssuer() string {
	return c.StandardClaims.Issuer
}

func (c *Claims) GetAudience() string {
	return c.StandardClaims.Audience
}

func (c *Claims) GetTokenId() string {
	return c.StandardClaims.Id
}
