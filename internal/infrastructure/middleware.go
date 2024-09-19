package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/higansama/xyz-multi-finance/internal/middleware"
)

type Middleware struct {
	AuthMiddleware      func(ctx *gin.Context)
	XAuthMiddleware     func(ctx *gin.Context)
	AuthOwnerMiddleware func(ctx *gin.Context)
}

func (infra *Infrastructure) setupMiddleware() (Middleware, error) {
	authMiddleware, err := middleware.NewAuthMiddleware(infra.Config)
	if err != nil {
		return Middleware{}, err
	}

	xAuthMiddleware, err := middleware.NewXAuthMiddleware(infra.Config)
	if err != nil {
		return Middleware{}, err
	}

	AuthOwerMiddleware, err := middleware.NewAuthOwnerMiddleware(infra.Config)
	if err != nil {
		return Middleware{}, err
	}

	return Middleware{
		AuthMiddleware:      authMiddleware.Handle,
		XAuthMiddleware:     xAuthMiddleware.Handle,
		AuthOwnerMiddleware: AuthOwerMiddleware.Handle,
	}, nil
}
