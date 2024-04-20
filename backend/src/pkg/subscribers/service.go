package subscribers

import (
	"newsletter/src/pkg/entity"
	"newsletter/src/pkg/utils/convert"
	newsletterError "newsletter/src/pkg/utils/error"
	"newsletter/src/pkg/utils/logger"
)

type UseCase interface {
	GetAllSubscribers() ([]entity.Subscribers, *newsletterError.ErrorCode)
	FindByEmail(email string) ([]entity.Subscribers, *newsletterError.ErrorCode)
	Insert(subscriber entity.Subscribers) *newsletterError.ErrorCode
	UpdateById(subscriber entity.Subscribers) *newsletterError.ErrorCode
}

type Service struct {
	UseCase
	Repo Repository
	Logs logger.Logger
}

func NewService(repo Repository, logs logger.Logger) *Service {
	service := &Service{
		Repo: repo,
		Logs: logs,
	}
	service.UseCase = service
	return service
}

func (service *Service) GetAllSubscribers() ([]entity.Subscribers, *newsletterError.ErrorCode) {
	res, err := service.Repo.GetAllSubscribers()

	if err != nil {
		return nil, convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
	}

	if len(res) == 0 {
		return nil, convert.ValueToErrorCodePointer(newsletterError.DataNotFound)
	}

	return res, nil
}

func (service *Service) FindByEmail(email string) ([]entity.Subscribers, *newsletterError.ErrorCode) {
	res, err := service.Repo.FindByEmail(email)

	if err != nil {
		return nil, convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
	}

	if len(res) == 0 {
		return nil, convert.ValueToErrorCodePointer(newsletterError.DataNotFound)
	}

	return res, nil
}

func (service *Service) Insert(subscriber entity.Subscribers) *newsletterError.ErrorCode {
	err := service.Repo.Insert(subscriber)

	if err != nil {
		return convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
	}

	return nil
}

func (service *Service) UpdateById(subscriber entity.Subscribers) *newsletterError.ErrorCode {
	err := service.Repo.UpdateById(subscriber)

	if err != nil {
		return convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
	}

	return nil
}
