package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/higansama/xyz-multi-finance/config"
	"github.com/higansama/xyz-multi-finance/internal/auth"
	"github.com/pkg/errors"
)

type AuthMiddleware struct {
	config config.Config
}

func NewAuthMiddleware(cfg config.Config) (AuthMiddleware, error) {
	return AuthMiddleware{config: cfg}, nil
}

func (m *AuthMiddleware) Handle(ctx *gin.Context) {
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

			ctx.Set("AUTH_DATA", *claims)

			ctx.Next()
			return
		}
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": "Unauthorized",
	})
	return
}

// CreateJwtToken For testing purpose
func (m *AuthMiddleware) CreateJwtToken() (string, error) {
	claims := &auth.JwtClaims{
		UserId:           "5df70a27127bfc211ceb1a46",
		UserName:         "",
		UserEmail:        "",
		UserRole:         "",
		RegisteredClaims: jwt.RegisteredClaims{ID: uuid.New().String(), ExpiresAt: jwt.NewNumericDate(carbon.Now().AddYears(60).ToStdTime())},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwToken, err := token.SignedString([]byte(m.config.Auth.JwtSecret))
	if err != nil {
		return jwToken, err
	}

	return jwToken, nil
}
