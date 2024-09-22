package auth

import (
	"github.com/gin-gonic/gin"
)

func GetAuthData(ctx *gin.Context) JwtClaims {
	val, _ := ctx.Get("AUTH_DATA")

	if claims, ok := val.(JwtClaims); ok {
		claims.UserId = claims.UserId
		claims.Salary = claims.Salary
		claims.UserEmail = claims.UserEmail
		claims.UserName = claims.UserName
		// claims.UserIsAMaster = claims.OldUserIsAMaster
		return claims
	}

	return JwtClaims{}
}
