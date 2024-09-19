// Code taken from official mongo package for Go
// For license and copyright see go.mongodb.org/mongo-driver

package mongodb

import "context"

type backgroundContext struct {
	context.Context
	childValuesCtx context.Context
}

func NewBackgroundContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return &backgroundContext{
		Context:        context.Background(),
		childValuesCtx: ctx,
	}
}

func (b *backgroundContext) Value(key any) any {
	return b.childValuesCtx.Value(key)
}
