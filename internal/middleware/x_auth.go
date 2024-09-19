package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/higansama/xyz-multi-finance/config"
)

type XAuthMiddleware struct {
	config config.Config
}

func NewXAuthMiddleware(cfg config.Config) (XAuthMiddleware, error) {
	return XAuthMiddleware{config: cfg}, nil
}

func (m *XAuthMiddleware) Handle(ctx *gin.Context) {
	authId := ctx.GetHeader("x-auth-id")
	authUserType := ctx.GetHeader("x-auth-user-type")
	aum := ctx.GetHeader("x-auth-user-master")

	if authId != "" && authUserType != "" && aum != "" {
		// authUserMaster, _ := strconv.ParseBool(aum)
		// authData := auth.JwtClaims{
		// 	UserId:           ,
		// 	UserName:         "",
		// 	UserEmail:        "",
		// 	UserRole:         "",
		// 	RegisteredClaims: jwt.RegisteredClaims{ID: authId},
		// }

		// ctx.Set("AUTH_DATA", authData)
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": "Unauthorized",
	})
	return
}
