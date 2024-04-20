package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	requestHeader "newsletter/src/api/requestheader"
	"newsletter/src/api/subscribers/handler"
	"newsletter/src/pkg/entity"
	"newsletter/src/pkg/subscribers/mocks"
	"newsletter/src/pkg/utils/convert"
	newsletterError "newsletter/src/pkg/utils/error"
	loggerMocks "newsletter/src/pkg/utils/logger/mocks"
	"newsletter/src/pkg/utils/mocker"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	uri            string
	service        *mocks.UseCase
	carTypeHandler *handler.SubscribersHandler
	logs           *loggerMocks.Logger
	recorder       *httptest.ResponseRecorder
	request        *http.Request
	router         *mux.Router

	mockServiceGetAllSubscribers *mocker.MockCall
)

func callServiceGetAllSubscribers() *mock.Call {
	return service.On("GetAllSubscribers")
}

func beforeEach() {
	uri = "/subscribers"
	service = &mocks.UseCase{}
	logs = &loggerMocks.Logger{}

	logs.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	carTypeHandler = &handler.SubscribersHandler{
		Service: service,
		Logs:    logs,
	}
}

func TestHandler_MakeSubscribersHandler(t *testing.T) {

	t.Run("should return struct subscribers handler when call make subscribers handler", func(t *testing.T) {
		beforeEach()
		handlerParam := handler.HandlerParam{
			Service: service,
			Logs:    logs,
		}

		handlerMakeSubscribers := handler.MakeSubscribersHandler(handlerParam)

		expectedResult := &handler.SubscribersHandler{
			Service: handlerParam.Service,
			Logs:    handlerParam.Logs,
		}
		assert.Equal(t, expectedResult, handlerMakeSubscribers)
	})
}

func TestHandler_GetAllSubscribers(t *testing.T) {
	beforeEachGetAllSubscribers := func() {
		beforeEach()
		router = mux.NewRouter()
		router.HandleFunc(uri, carTypeHandler.GetAllSubscribers)
		recorder = httptest.NewRecorder()
		request = httptest.NewRequest(http.MethodGet, uri, nil)

		mockServiceGetAllSubscribers = mocker.NewMockCall(callServiceGetAllSubscribers)
		mockServiceGetAllSubscribers.Return(nil, nil)
	}

	t.Run("should call service search when request get all subscribers", func(t *testing.T) {
		beforeEachGetAllSubscribers()
		request = httptest.NewRequest(http.MethodGet, uri, nil)

		router.ServeHTTP(recorder, request)

		service.AssertCalled(t, "GetAllSubscribers")
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})

	t.Run("should response internal server error when service get all subscribers failed", func(t *testing.T) {
		beforeEachGetAllSubscribers()
		err := newsletterError.InternalServerError
		mockServiceGetAllSubscribers.Return([]entity.Subscribers{}, &err)

		router.ServeHTTP(recorder, request)

		expectedStatusCode, expectedError := newsletterError.MapMessageError(newsletterError.InternalServerError, "en")
		var body newsletterError.Error
		json.NewDecoder(recorder.Body).Decode(&body)
		assert.Equal(t, expectedError, body)
		assert.Equal(t, expectedStatusCode, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})

	t.Run("should response data not found error when service get all subscribers failed", func(t *testing.T) {
		beforeEachGetAllSubscribers()
		mockServiceGetAllSubscribers.Return(nil, convert.ValueToErrorCodePointer(newsletterError.DataNotFound))

		router.ServeHTTP(recorder, request)

		expectedStatusCode, expectedError := newsletterError.MapMessageError(newsletterError.DataNotFound, "en")
		var body newsletterError.Error
		json.NewDecoder(recorder.Body).Decode(&body)
		assert.Equal(t, expectedError, body)
		assert.Equal(t, expectedStatusCode, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})

	t.Run("should response get all subscribers when request service get all subscribers success", func(t *testing.T) {
		beforeEachGetAllSubscribers()
		resSubscribers := []entity.Subscribers{
			{
				ID:   1,
				Name: "Test1",
			},
			{
				ID:   2,
				Name: "Test1",
			},
		}
		mockServiceGetAllSubscribers.Return(resSubscribers, nil)

		router.ServeHTTP(recorder, request)

		expectedBody := resSubscribers

		var body []entity.Subscribers
		json.NewDecoder(recorder.Body).Decode(&body)
		assert.Equal(t, expectedBody, body)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})
}
