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
	mockRepoInsert            *mocker.MockCall
	mockRepoUpdateByEmail     *mocker.MockCall
)

func callRepoGetAllSubscribers() *mock.Call {
	return repository.On("GetAllSubscribers")
}

func callRepoFindByEmail() *mock.Call {
	return repository.On("FindByEmail", mock.Anything)
}

func callRepoInsert() *mock.Call {
	return repository.On("Insert", mock.Anything)
}

func callRepoUpdateByEmail() *mock.Call {
	return repository.On("UpdateByEmail", mock.Anything)
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

func TestService_FindByEmail(t *testing.T) {
	beforeEachFindByEmail := func() {
		beforeEach()

		mockRepoFindByEmail = mocker.NewMockCall(callRepoFindByEmail)
		mockRepoFindByEmail.Return(nil, nil)
	}

	t.Run("should call repository find by email when call service find by email", func(t *testing.T) {
		beforeEachFindByEmail()

		email := "ajistestmail@gmail.com"

		service.FindByEmail(email)

		repository.AssertCalled(t, "FindByEmail", email)
	})

	t.Run("should response internal server error when repository find by email failed", func(t *testing.T) {
		beforeEachFindByEmail()
		mockRepoFindByEmail.Return(nil, errors.New("Error"))

		email := "ajistestmail@gmail.com"
		res, err := service.FindByEmail(email)

		expectedError := convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, res)
	})

	t.Run("should return data not found when data not founded", func(t *testing.T) {
		beforeEachFindByEmail()
		mockRepoFindByEmail.Return(nil, nil)

		email := "ajistestmail@gmail.com"
		res, err := service.FindByEmail(email)

		expectedError := convert.ValueToErrorCodePointer(newsletterError.DataNotFound)
		assert.Nil(t, res)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return data when found data", func(t *testing.T) {
		beforeEachFindByEmail()
		mockSubscribers := []entity.Subscribers{
			{
				ID:    1,
				Name:  "TEST",
				Email: "ajistestmail@gmail.com",
			},
		}
		mockRepoFindByEmail.Return(mockSubscribers, nil)

		email := "ajistestmail@gmail.com"
		res, err := service.FindByEmail(email)

		assert.Equal(t, mockSubscribers, res)
		assert.Nil(t, err)
	})

}

func TestService_Insert(t *testing.T) {
	beforeEachInsert := func() {
		beforeEach()

		mockRepoInsert = mocker.NewMockCall(callRepoInsert)
		mockRepoInsert.Return(nil)
	}

	t.Run("should call repository insert when call service insert", func(t *testing.T) {
		beforeEachInsert()
		mockSubscribers := entity.Subscribers{
			Name: "test",
		}

		service.Insert(mockSubscribers)

		repository.AssertCalled(t, "Insert", mockSubscribers)
	})

	t.Run("should response internal server error when repository insert failed", func(t *testing.T) {
		beforeEachInsert()
		mockSubscribers := entity.Subscribers{
			Name:  "test",
			Email: "ajistestmail@gmail.com",
		}
		mockRepoInsert.Return(errors.New("Error"))

		err := service.Insert(mockSubscribers)

		expectedError := convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return error nil when call repository insert success", func(t *testing.T) {
		beforeEachInsert()
		mockSubscribers := entity.Subscribers{
			Name:  "test",
			Email: "ajistestmail@gmail.com",
		}

		err := service.Insert(mockSubscribers)

		assert.Nil(t, err)
	})
}

func TestService_UpdateByEmail(t *testing.T) {
	beforeEachUpdateByEmail := func() {
		beforeEach()

		mockRepoUpdateByEmail = mocker.NewMockCall(callRepoUpdateByEmail)
		mockRepoUpdateByEmail.Return(nil)
	}

	t.Run("should call repository update by id when call service update by id", func(t *testing.T) {
		beforeEachUpdateByEmail()
		mockSubscribers := entity.Subscribers{
			Name:  "test",
			Email: "ajistestmail@gmail.com",
		}

		service.UpdateByEmail(mockSubscribers)

		repository.AssertCalled(t, "UpdateByEmail", mockSubscribers)
	})

	t.Run("should return internal server error when call repository update by id failed", func(t *testing.T) {
		beforeEachUpdateByEmail()
		mockSubscribers := entity.Subscribers{
			Name:  "test",
			Email: "ajistestmail@gmail.com",
		}
		mockRepoUpdateByEmail.Return(errors.New("Error"))

		err := service.UpdateByEmail(mockSubscribers)

		expectedError := convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return error nil when call repository update by id success", func(t *testing.T) {
		beforeEachUpdateByEmail()
		mockSubscribers := entity.Subscribers{
			Name:  "test",
			Email: "ajistestmail@gmail.com",
		}

		err := service.UpdateByEmail(mockSubscribers)

		assert.Nil(t, err)
	})
}
