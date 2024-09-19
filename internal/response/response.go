package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/higansama/xyz-multi-finance/config"
	ierrors "github.com/higansama/xyz-multi-finance/internal/errors"
	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

func CreatedResponse(ctx *gin.Context, msg string, m any) {
	ctx.JSON(http.StatusCreated, gin.H{
		"message": utils.UpperFirst(msg),
		"data":    m,
	})
}

func CreatedResponseMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusCreated, gin.H{
		"message": utils.UpperFirst(msg),
	})
}

func OkResponseData(ctx *gin.Context, m any) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": m,
	})
}

func OkResponseWith(ctx *gin.Context, msg string, m any) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": utils.UpperFirst(msg),
		"data":    m,
	})
}

func OkResponse(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": utils.UpperFirst(msg),
	})
}

func BadRequestResponse(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": utils.UpperFirst(msg),
	})

	defer ctx.AbortWithStatus(http.StatusBadRequest)
}

func NewResponseError(ctx *gin.Context, err error, hecm ...map[error]int) {
	if err == nil {
		panic("Error cannot be nil")
	}
	httpStatusCode := http.StatusInternalServerError
	errorCode := ""

	hscOv := -1
	httpErrCodeMap := make(map[error]int)
	if len(hecm) > 0 {
		httpErrCodeMap = hecm[0]
	}
	for k, v := range httpErrCodeMap {
		if errors.Is(err, k) {
			hscOv = v
			break
		}
	}

	err = ierrors.DomainErrorToResponseError(err)
	errMessage := err.Error()

	var vErrs ierrors.ValidationError
	var respErr ierrors.ResponseError

	if errors.As(err, &vErrs) {
		ValidationErrorResponse(ctx, vErrs)
		return
	} else if errors.As(err, &respErr) {
		httpStatusCode = respErr.Code
		errorCode = respErr.ErrCode
		errMessage = respErr.Message
		if respErr.Err != nil {
			err = respErr.Err // set original error back
		}
	} else if errors.Is(err, ierrors.ErrInternalServerError) {
		httpStatusCode = http.StatusInternalServerError
	}

	if hscOv != -1 { // http status code override
		httpStatusCode = hscOv
	}
	if httpStatusCode == http.StatusInternalServerError {
		log.Error().Stack().Err(err).Send()

		if !config.Cfg.App.Debug {
			err = ierrors.ErrInternalServerError
		}
	}

	json := gin.H{
		"message": utils.UpperFirst(errMessage),
	}
	if errorCode != "" {
		json["error_code"] = errorCode
	}

	ctx.JSON(httpStatusCode, json)

	defer ctx.AbortWithStatus(httpStatusCode)
}

func ValidationErrorResponse(ctx *gin.Context, v ierrors.ValidationError) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": "The given data was invalid",
		"errors":  v.Errors,
	})

	defer ctx.AbortWithStatus(http.StatusUnprocessableEntity)
}

func RenderExcelize(ctx *gin.Context, file *excelize.File, filename string) {
	ctx.Status(http.StatusOK)

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	file.Write(ctx.Writer)
}
