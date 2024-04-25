package logger

type Logger interface {
	Error(endpoint, actionName string, request, message interface{})
}
