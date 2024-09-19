package validator

import (
	"mime/multipart"

	ierrors "github.com/higansama/xyz-multi-finance/internal/errors"
	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/pkg/errors"
)

var contentTypeAlias = map[string]string{
	"gif":              "image/gif",
	"jpeg":             "image/jpeg",
	"png":              "image/png",
	"csv":              "text/csv",
	"plain":            "text/plain",
	"x-zip":            "multipart/x-zip",
	"x-zip-compressed": "application/x-zip-compressed",
	"x-7z-compressed":  "application/x-7z-compressed",
	"x-rar-compressed": "application/x-rar-compressed",
	"zip":              "application/zip",
	"gzip":             "application/gzip",
	"pdf":              "application/pdf",
	"doc":              "application/msword",
	"xls":              "application/vnd.ms-excel",
	"ppt":              "application/vnd.ms-powerpoint",
	"docx":             "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"pptx":             "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"xlsx":             "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
}
var contentTypeMap = map[string]bool{
	"video/wav":       true,
	"video/mpeg":      true,
	"video/quicktime": true,
	"video/mp4":       true,
	"video/ogg":       true,
	"video/webm":      true,
	"audio/webm":      true,
	"audio/ogg":       true,
	"audio/mp3":       true,
	"audio/mp4":       true,
	"audio/mpeg":      true,
}

func init() {
	for _, v := range contentTypeAlias {
		contentTypeMap[v] = true
	}
}

func ValidateFileContentType(file *multipart.FileHeader, contentTypes []string) (bool, error) {
	if file == nil {
		return false, errors.New("file is not exists")
	}

	ctMap := make(map[string]bool)
	for _, ct := range contentTypes {
		if v, ok := contentTypeAlias[ct]; ok { // convert alias to actual
			ct = v
		}
		if !contentTypeMap[ct] {
			return false, errors.Errorf("unknown '%s' content type", ct)
		}
		ctMap[ct] = true
	}

	ct := utils.GetFileContentType(file)
	if !ctMap[ct] {
		return false, nil
	}

	return true, nil
}

func ValidateBool(val string) (string, *ierrors.FieldError) {
	validBool := map[string]bool{
		"true":  true,
		"false": true,
		"1":     true,
		"0":     true,
	}
	if !validBool[val] {
		return val, &ierrors.FieldError{
			Msg: "Must be a boolean",
			Tag: "bool",
		}
	}

	if val == "1" {
		val = "true"
	} else if val == "0" {
		val = "false"
	}

	return val, nil
}
