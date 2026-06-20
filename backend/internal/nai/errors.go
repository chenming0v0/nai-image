package nai

// ValidationError 表示请求参数校验失败，附带 HTTP 状态码建议。
type ValidationError struct {
	Status int
	Field  string
	Msg    string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Msg
	}
	return e.Msg
}

func newValidationError(field, msg string) *ValidationError {
	return &ValidationError{Status: 400, Field: field, Msg: msg}
}

// 上游错误，携带 HTTP 状态码与原始信息。
type UpstreamError struct {
	Status  int
	Message string
}

func (e *UpstreamError) Error() string {
	return e.Message
}
