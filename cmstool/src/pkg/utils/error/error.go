package error

type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewError(code ErrorCode, message string, data ...interface{}) *Error {
	if data != nil {
		return &Error{Code: code, Message: message, Data: data[0]}
	}

	return &Error{Code: code, Message: message}
}
