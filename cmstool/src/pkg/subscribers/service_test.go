package subscribers_test

import (
	"errors"
	"subscribetool/src/pkg/entity"
	subscribers "subscribetool/src/pkg/subscribers"
	"subscribetool/src/pkg/subscribers/mocks"
	"subscribetool/src/pkg/utils/convert"
	"subscribetool/src/pkg/utils/email"
	emailMocks "subscribetool/src/pkg/utils/email/mocks"
	subscribetoolError "subscribetool/src/pkg/utils/error"
	loggerMocks "subscribetool/src/pkg/utils/logger/mocks"
	"subscribetool/src/pkg/utils/mocker"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockUseCase       *mocks.UseCase
	repository        *mocks.Repository
	service           *subscribers.Service
	logs              *loggerMocks.Logger
	utilsEmailService *emailMocks.UseCase

	mockRepoGetAllSubscribers    *mocker.MockCall
	mockServiceGetAllSubscribers *mocker.MockCall
	mockUtilsEmailServiceSend    *mocker.MockCall
)

func callRepoGetAllSubscribers() *mock.Call {
	return repository.On("GetAllSubscribers")
}

func callServiceGetAllSubscribers() *mock.Call {
	return mockUseCase.On("GetAllSubscribers")
}

func callUtilsEmailServiceSend() *mock.Call {
	return utilsEmailService.On("Send", mock.Anything)
}

func beforeEach() {
	mockUseCase = &mocks.UseCase{}
	repository = &mocks.Repository{}
	logs = &loggerMocks.Logger{}
	utilsEmailService = &emailMocks.UseCase{}

	logs.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	service = &subscribers.Service{
		Repo:              repository,
		Logs:              logs,
		UtilsEmailService: utilsEmailService,
	}
	service.UseCase = mockUseCase
}

func TestService_NewService(t *testing.T) {
	t.Run("should return struct subscribers service when call new service", func(t *testing.T) {
		beforeEach()
		serviceParam := subscribers.ServiceParam{
			Repo:              repository,
			Logs:              logs,
			UtilsEmailService: utilsEmailService,
		}
		resService := subscribers.NewService(serviceParam)

		expectedService := &subscribers.Service{
			Repo:              repository,
			Logs:              logs,
			UtilsEmailService: utilsEmailService,
		}
		expectedService.UseCase = expectedService

		assert.Equal(t, expectedService, resService)
	})
}

func TestService_GetAllSubscribers(t *testing.T) {
	beforeEachGetAllSubscribers := func() {
		beforeEach()

		mockRepoGetAllSubscribers = mocker.NewMockCall(callRepoGetAllSubscribers)
		mockRepoGetAllSubscribers.Return(nil, nil)
	}

	t.Run("should call repository get all subscribers when call service get all subscribers", func(t *testing.T) {
		beforeEachGetAllSubscribers()

		service.GetAllSubscribers()

		repository.AssertCalled(t, "GetAllSubscribers")
	})

	t.Run("should response internal server error when repository get all subscribers failed", func(t *testing.T) {
		beforeEachGetAllSubscribers()
		mockRepoGetAllSubscribers.Return(nil, errors.New("Error"))

		res, err := service.GetAllSubscribers()

		expectedError := convert.ValueToErrorCodePointer(subscribetoolError.InternalServerError)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, res)
	})

	t.Run("should return data not found when data not founded", func(t *testing.T) {
		beforeEachGetAllSubscribers()
		mockRepoGetAllSubscribers.Return(nil, nil)

		res, err := service.GetAllSubscribers()

		expectedError := convert.ValueToErrorCodePointer(subscribetoolError.DataNotFound)
		assert.Nil(t, res)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return data when found data", func(t *testing.T) {
		beforeEachGetAllSubscribers()
		mockSubscribers := []entity.Subscribers{
			{
				ID:   1,
				Name: "TEST",
			},
		}
		mockRepoGetAllSubscribers.Return(mockSubscribers, nil)

		res, err := service.GetAllSubscribers()

		assert.Equal(t, mockSubscribers, res)
		assert.Nil(t, err)
	})

}

func TestService_SentEmail(t *testing.T) {
	beforeEachSentEmail := func() {
		beforeEach()

		mockServiceGetAllSubscribers = mocker.NewMockCall(callServiceGetAllSubscribers)
		mockServiceGetAllSubscribers.Return(nil, nil)
		mockUtilsEmailServiceSend = mocker.NewMockCall(callUtilsEmailServiceSend)
		mockUtilsEmailServiceSend.Return(nil)

	}

	t.Run("should call service get all subscribers when call service sent email", func(t *testing.T) {
		beforeEachSentEmail()

		service.SentEmail()

		mockUseCase.AssertCalled(t, "GetAllSubscribers")
	})

	t.Run("should return internal server error when call service get all subscribers failed", func(t *testing.T) {
		beforeEachSentEmail()
		mockServiceGetAllSubscribers.Return(nil, convert.ValueToErrorCodePointer(subscribetoolError.InternalServerError))

		err := service.SentEmail()

		expectedError := convert.ValueToErrorCodePointer(subscribetoolError.InternalServerError)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return data not found when call service get all subscribers response zero lenght", func(t *testing.T) {
		beforeEachSentEmail()
		mockServiceGetAllSubscribers.Return([]entity.Subscribers{}, nil)

		err := service.SentEmail()

		expectedError := convert.ValueToErrorCodePointer(subscribetoolError.DataNotFound)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should call utils email service send when call service sent email", func(t *testing.T) {
		beforeEachSentEmail()

		mockDataSubscribers := mockDataSubscribers()
		mockServiceGetAllSubscribers.Return(mockDataSubscribers, nil)

		service.SentEmail()

		emailTarget := []string{mockDataSubscribers[0].Email}

		mailInfo := email.SentMailContent{
			To:      emailTarget,
			Supject: "Test sent mail for subscribers",
			Body:    "This email sent for notification. Test sent mail for subscribers",
		}

		utilsEmailService.AssertCalled(t, "Send", mailInfo)
	})

	t.Run("should return internal server error when call utils email service send failed", func(t *testing.T) {
		beforeEachSentEmail()

		mockDataSubscribers := mockDataSubscribers()
		mockServiceGetAllSubscribers.Return(mockDataSubscribers, nil)
		mockUtilsEmailServiceSend.Return(errors.New("error"))

		err := service.SentEmail()

		expectedError := convert.ValueToErrorCodePointer(subscribetoolError.InternalServerError)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return error nil when call service sent mail success", func(t *testing.T) {
		beforeEachSentEmail()
		mockDataSubscribers := mockDataSubscribers()
		mockServiceGetAllSubscribers.Return(mockDataSubscribers, nil)
		mockUtilsEmailServiceSend.Return(nil)

		err := service.SentEmail()

		assert.Nil(t, err)
	})

}

func mockDataSubscribers() []entity.Subscribers {
	return []entity.Subscribers{
		{
			ID:           1,
			Name:         "test",
			Email:        "ajistestmail@gmail.com",
			IsSubscribed: true,
		},
	}
}
