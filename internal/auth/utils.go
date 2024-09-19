package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetAuthData(ctx *gin.Context) JwtClaims {
	val, _ := ctx.Get("AUTH_DATA")

	if claims, ok := val.(JwtClaims); ok {
		fmt.Println(claims)
		claims.UserId = claims.UserId
		claims.UserRole = claims.UserRole
		claims.UserEmail = claims.UserEmail
		claims.UserName = claims.UserName
		// claims.UserIsAMaster = claims.OldUserIsAMaster
		return claims
	}

	return JwtClaims{}
}
