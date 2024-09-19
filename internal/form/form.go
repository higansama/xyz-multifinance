package form

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strings"
)

type File struct {
	Content     *multipart.FileHeader
	ContentType string
	Unfilled    bool // unfilled means the field exist, but no value provided
}

func GetTrimmedPostForm(ctx *gin.Context, field string) string {
	val := ctx.PostForm(field)

	return strings.TrimSpace(val)
}

func GetNullableStringPostForm(ctx *gin.Context, field string, trim bool) *string {
	val, ok := ctx.GetPostForm(field)
	if trim {
		val = GetTrimmedPostForm(ctx, field)
	}
	res := &val
	if !ok {
		res = nil
	}
	return res
}
