package subscribers_test

import (
	"errors"
	"newsletter/src/pkg/entity"
	subscribers "newsletter/src/pkg/subscribers"
	"newsletter/src/pkg/subscribers/mocks"
	"newsletter/src/pkg/utils/convert"
	newsletterError "newsletter/src/pkg/utils/error"
	loggerMocks "newsletter/src/pkg/utils/logger/mocks"
	"newsletter/src/pkg/utils/mocker"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockUseCase *mocks.UseCase
	repository  *mocks.Repository
	service     *subscribers.Service
	logs        *loggerMocks.Logger

	mockRepoGetAllSubscribers *mocker.MockCall
	mockRepoFindByEmail       *mocker.MockCall
)

func callRepoGetAllSubscribers() *mock.Call {
	return repository.On("GetAllSubscribers")
}

func callRepoFindByEmail() *mock.Call {
	return repository.On("FindByEmail")
}

func beforeEach() {
	mockUseCase = &mocks.UseCase{}
	repository = &mocks.Repository{}
	logs = &loggerMocks.Logger{}

	logs.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	service = &subscribers.Service{
		Repo: repository,
		Logs: logs,
	}
	service.UseCase = mockUseCase
}

func TestService_NewService(t *testing.T) {
	t.Run("should return struct subscribers service when call new service", func(t *testing.T) {
		beforeEach()
		resService := subscribers.NewService(repository, logs)

		expectedService := &subscribers.Service{
			Repo: repository,
			Logs: logs,
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

		expectedError := convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, res)
	})

	t.Run("should return data not found when data not founded", func(t *testing.T) {
		beforeEachGetAllSubscribers()
		mockRepoGetAllSubscribers.Return(nil, nil)

		res, err := service.GetAllSubscribers()

		expectedError := convert.ValueToErrorCodePointer(newsletterError.DataNotFound)
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
