package request

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ierrors "github.com/higansama/xyz-multi-finance/internal/errors"
	"github.com/higansama/xyz-multi-finance/internal/processor"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Request struct {
	structProcessor processor.StructProcessor
}

func NewRequest(structProcessor processor.StructProcessor) Request {
	return Request{
		structProcessor: structProcessor,
	}
}

type BindOptions struct {
	ContentType string
	ShouldLog   bool
}

func (r *Request) Bind(ctx *gin.Context, obj any, options ...BindOptions) error {
	var opts BindOptions
	if len(options) > 0 {
		opts = options[0]
	}

	contentType := ctx.ContentType()
	if opts.ContentType != "" {
		contentType = opts.ContentType
	}

	b := binding.Default(ctx.Request.Method, contentType)

	fn := func() error {
		validator := binding.Validator
		// temporary unset validator, since we don't want validator runs when bind the data
		binding.Validator = nil
		err := r.checkBindingError(ctx.ShouldBindWith(obj, b))
		if err != nil {
			return err
		}

		if ctx.Request.Method == http.MethodGet && contentType == binding.MIMEJSON {
			err := r.checkBindingError(ctx.ShouldBindJSON(obj))
			if err != nil {
				return err
			}
		}

		err = r.checkBindingError(ctx.ShouldBindUri(obj))
		if err != nil {
			return err
		}
		binding.Validator = validator

		err = r.structProcessor.ProcessAndValidate(ctx, obj)
		if err != nil {
			return err
		}

		return nil
	}

	err := fn()
	if err != nil {
		if opts.ShouldLog {
			var vErrs ierrors.ValidationError
			if errors.As(err, &vErrs) {
				log.Error().Any("Errors", vErrs.Errors).Msg("Validation error")
			} else {
				log.Error().Err(err).Msg("Binding error")
			}
		}
		return err
	}

	return nil
}

func (r *Request) checkBindingError(err error) error {
	if err != nil {
		if errors.Is(err, http.ErrMissingBoundary) ||
			errors.Is(err, http.ErrNotMultipart) ||
			errors.Is(err, io.EOF) ||
			err.Error() == "unknown type" {
			// ignore this type of error
		} else {
			return errors.WithStack(ierrors.BadRequest("the given data was invalid"))
		}
	}
	return nil
}

// Log is used when we need to debug the raw request payload.
// Please use carefully, since it's using io.ReadAll.
func (r *Request) Log(ctx *gin.Context, msg string) {
	body, _ := io.ReadAll(ctx.Request.Body)

	log.Debug().
		Str("Plain Request Payload", string(body)).
		Msg(msg)

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
}
