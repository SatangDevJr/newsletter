package handler

import (
	"encoding/json"
	"net/http"
	requestHeader "newsletter/src/api/requestheader"
	subscribers "newsletter/src/pkg/subscribers"
	newsletterError "newsletter/src/pkg/utils/error"
	"newsletter/src/pkg/utils/logger"
)

type SubscribersHandler struct {
	Service subscribers.UseCase
	Logs    logger.Logger
}

func MakeSubscribersHandler(handlerParam HandlerParam) *SubscribersHandler {
	return &SubscribersHandler{
		Service: handlerParam.Service,
		Logs:    handlerParam.Logs,
	}
}

func (handler *SubscribersHandler) GetAllSubscribers(response http.ResponseWriter, request *http.Request) {

	response.Header().Set(requestHeader.ContentType, requestHeader.ApplicationJson)

	res, err := handler.Service.GetAllSubscribers()
	if err != nil {
		switch *err {
		case newsletterError.DataNotFound:
			go handler.Logs.Error(request.URL.Path, "subscribers_handler_getAllSubscribers_DataNotFound", nil, err)
		default:
			go handler.Logs.Error(request.URL.Path, "subscribers_handler_getAllSubscribers_InternalServerError", nil, err)
		}
		statusCode, errMsg := newsletterError.MapMessageError(*err, "en")
		response.WriteHeader(statusCode)
		json.NewEncoder(response).Encode(&errMsg)
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(&res)
}
