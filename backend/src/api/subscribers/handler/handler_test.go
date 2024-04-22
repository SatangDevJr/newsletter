package handler_test

import (
	"bytes"
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
	uri               string
	service           *mocks.UseCase
	subscriberHandler *handler.SubscribersHandler
	logs              *loggerMocks.Logger
	recorder          *httptest.ResponseRecorder
	request           *http.Request
	router            *mux.Router

	mockServiceGetAllSubscribers *mocker.MockCall
	mockServiceSubscribe         *mocker.MockCall
	mockServiceUnsubscribe       *mocker.MockCall
)

func callServiceGetAllSubscribers() *mock.Call {
	return service.On("GetAllSubscribers")
}

func callServiceSubscribe() *mock.Call {
	return service.On("Subscribe", mock.Anything)
}

func callServiceUnsubscribe() *mock.Call {
	return service.On("Unsubscribe", mock.Anything)
}

func beforeEach() {
	uri = "/subscribers"
	service = &mocks.UseCase{}
	logs = &loggerMocks.Logger{}

	logs.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	subscriberHandler = &handler.SubscribersHandler{
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
		router.HandleFunc(uri, subscriberHandler.GetAllSubscribers)
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

func TestHandler_Subscribe(t *testing.T) {
	beforeEachSubscribe := func() {
		beforeEach()
		router = mux.NewRouter()
		router.HandleFunc(uri, subscriberHandler.Subscribe)
		recorder = httptest.NewRecorder()
		request = httptest.NewRequest(http.MethodGet, uri, nil)

		mockServiceSubscribe = mocker.NewMockCall(callServiceSubscribe)
		mockServiceSubscribe.Return(nil, nil)
	}

	t.Run("should response bad request when request handler subscribe format", func(t *testing.T) {
		beforeEachSubscribe()
		requestBody := ``
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(requestBody)))

		router.ServeHTTP(recorder, request)

		var responseBody newsletterError.Error
		json.NewDecoder(recorder.Body).Decode(&responseBody)
		err := newsletterError.NewError(newsletterError.BadRequest, newsletterError.BadRequestMessage)
		assert.Equal(t, *err, responseBody)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("should call service subscribe when request subscribe", func(t *testing.T) {
		beforeEachSubscribe()

		subscribers := entity.Subscribers{
			Name:  "TEST",
			Email: "ajistestmail@gmail.com",
		}

		jsonMock, _ := json.Marshal(subscribers)
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(jsonMock)))

		router.ServeHTTP(recorder, request)

		service.AssertCalled(t, "Subscribe", subscribers)
	})

	t.Run("should response internal server error when service subscribe failed", func(t *testing.T) {
		beforeEachSubscribe()

		subscribers := entity.Subscribers{
			Name:  "TEST",
			Email: "ajistestmail@gmail.com",
		}

		jsonMock, _ := json.Marshal(subscribers)
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(jsonMock)))
		mockServiceSubscribe.Return(convert.ValueToErrorCodePointer(newsletterError.InternalServerError))

		router.ServeHTTP(recorder, request)

		expectedStatusCode, expectedError := newsletterError.MapMessageError(newsletterError.InternalServerError, "en")
		var responseBody newsletterError.Error
		json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.Equal(t, expectedError, responseBody)
		assert.Equal(t, expectedStatusCode, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))

	})

	t.Run("should response data not found server error when service subscribe failed", func(t *testing.T) {
		beforeEachSubscribe()

		subscribers := entity.Subscribers{
			Name:  "TEST",
			Email: "ajistestmail@gmail.com",
		}

		jsonMock, _ := json.Marshal(subscribers)
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(jsonMock)))
		mockServiceSubscribe.Return(convert.ValueToErrorCodePointer(newsletterError.DataNotFound))

		router.ServeHTTP(recorder, request)

		expectedStatusCode, expectedError := newsletterError.MapMessageError(newsletterError.DataNotFound, "en")
		var responseBody newsletterError.Error
		json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.Equal(t, expectedError, responseBody)
		assert.Equal(t, expectedStatusCode, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))

	})

	t.Run("should return ok when request subscribe success", func(t *testing.T) {
		beforeEachSubscribe()
		subscribers := entity.Subscribers{
			Name:  "TEST",
			Email: "ajistestmail@gmail.com",
		}
		jsonMock, _ := json.Marshal(subscribers)
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(jsonMock)))

		router.ServeHTTP(recorder, request)

		res := handler.ResponseSucess{
			Body: "subscribe success",
		}

		var responseBody handler.ResponseSucess
		json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.Equal(t, res, responseBody)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})
}

func TestHandler_Unsubscribe(t *testing.T) {
	beforeEachSubscribe := func() {
		beforeEach()
		router = mux.NewRouter()
		router.HandleFunc(uri, subscriberHandler.Unsubscribe)
		recorder = httptest.NewRecorder()
		request = httptest.NewRequest(http.MethodGet, uri, nil)

		mockServiceUnsubscribe = mocker.NewMockCall(callServiceUnsubscribe)
		mockServiceUnsubscribe.Return(nil, nil)
	}

	t.Run("should response bad request when request handler unsubscribe format", func(t *testing.T) {
		beforeEachSubscribe()
		requestBody := ``
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(requestBody)))

		router.ServeHTTP(recorder, request)

		var responseBody newsletterError.Error
		json.NewDecoder(recorder.Body).Decode(&responseBody)
		err := newsletterError.NewError(newsletterError.BadRequest, newsletterError.BadRequestMessage)
		assert.Equal(t, *err, responseBody)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("should call service unsubscribe when request unsubscribe", func(t *testing.T) {
		beforeEachSubscribe()

		unsubscribers := entity.Subscribers{
			Name:  "TEST",
			Email: "ajistestmail@gmail.com",
		}

		jsonMock, _ := json.Marshal(unsubscribers)
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(jsonMock)))

		router.ServeHTTP(recorder, request)

		service.AssertCalled(t, "Unsubscribe", unsubscribers)
	})

	t.Run("should response internal server error when service unsubscribe failed", func(t *testing.T) {
		beforeEachSubscribe()

		unsubscribers := entity.Subscribers{
			Name:  "TEST",
			Email: "ajistestmail@gmail.com",
		}

		jsonMock, _ := json.Marshal(unsubscribers)
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(jsonMock)))
		mockServiceUnsubscribe.Return(convert.ValueToErrorCodePointer(newsletterError.InternalServerError))

		router.ServeHTTP(recorder, request)

		expectedStatusCode, expectedError := newsletterError.MapMessageError(newsletterError.InternalServerError, "en")
		var responseBody newsletterError.Error
		json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.Equal(t, expectedError, responseBody)
		assert.Equal(t, expectedStatusCode, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))

	})

	t.Run("should response data not found server error when service unsubscribe failed", func(t *testing.T) {
		beforeEachSubscribe()

		unsubscribers := entity.Subscribers{
			Name:  "TEST",
			Email: "ajistestmail@gmail.com",
		}

		jsonMock, _ := json.Marshal(unsubscribers)
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(jsonMock)))
		mockServiceUnsubscribe.Return(convert.ValueToErrorCodePointer(newsletterError.DataNotFound))

		router.ServeHTTP(recorder, request)

		expectedStatusCode, expectedError := newsletterError.MapMessageError(newsletterError.DataNotFound, "en")
		var responseBody newsletterError.Error
		json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.Equal(t, expectedError, responseBody)
		assert.Equal(t, expectedStatusCode, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))

	})

	t.Run("should return ok when request unsubscribe success", func(t *testing.T) {
		beforeEachSubscribe()
		unsubscribers := entity.Subscribers{
			Name:  "TEST",
			Email: "ajistestmail@gmail.com",
		}
		jsonMock, _ := json.Marshal(unsubscribers)
		request = httptest.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(jsonMock)))

		router.ServeHTTP(recorder, request)

		res := handler.ResponseSucess{
			Body: "unsubscribe success",
		}

		var responseBody handler.ResponseSucess
		json.NewDecoder(recorder.Body).Decode(&responseBody)
		assert.Equal(t, res, responseBody)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})
}
