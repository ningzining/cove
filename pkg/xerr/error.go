package xerr

import (
	"errors"
	"fmt"
)

type Error struct {
	Code         int
	TemplateDate map[string]interface{}
}

func New(code int) *Error {
	return &Error{Code: code}
}

// WithTemplateDate 设置错误的模板数据.
func (e *Error) WithTemplateDate(templateData map[string]interface{}) {
	e.TemplateDate = templateData
}

// Error 实现 error 接口中的 `Error` 方法.
func (e *Error) Error() string {
	return fmt.Sprintf("code: %d template: %v", e.Code, e.TemplateDate)
}

// Decode 尝试从 err 中解析出业务错误码和错误信息.
func Decode(err error) (int, map[string]interface{}) {
	if err == nil {
		return 200, nil
	}

	var xerr *Error
	if errors.As(err, &xerr) {
		return xerr.Code, xerr.TemplateDate
	}

	return 100000, nil
}
