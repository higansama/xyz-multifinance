package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/higansama/xyz-multi-finance/config"
	"github.com/higansama/xyz-multi-finance/internal/auth"
	"github.com/pkg/errors"
)

type AuthOwnerMiddleware struct {
	config config.Config
}

func NewAuthOwnerMiddleware(cfg config.Config) (AuthOwnerMiddleware, error) {
	return AuthOwnerMiddleware{config: cfg}, nil
}

func (m *AuthOwnerMiddleware) Handle(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		jwToken := parts[1]
		token, err := jwt.ParseWithClaims(
			jwToken,
			&auth.JwtClaims{},
			func(token *jwt.Token) (any, error) {
				if jwt.GetSigningMethod(jwt.SigningMethodHS256.Alg()) != token.Method {
					return nil, errors.Errorf("invalid signing method: %v", token.Header["alg"])
				}
				return []byte(m.config.Auth.JwtSecret), nil
			},
		)

		if err == nil {
			claims := token.Claims.(*auth.JwtClaims)
			// if claims.UserRole != "4" {
			// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			// 		"message": "Unauthorized",
			// 	})
			// }
			ctx.Set("AUTH_DATA", *claims)

			ctx.Next()
			return
		}
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": "Unauthorized",
	})
}
