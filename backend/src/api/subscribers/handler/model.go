package handler

import (
	subscribers "newsletter/src/pkg/subscribers"
	"newsletter/src/pkg/utils/logger"
)

type HandlerParam struct {
	Service subscribers.UseCase
	Logs    logger.Logger
}

type ResponseSucess struct {
	Body string `json:"body"`
}
