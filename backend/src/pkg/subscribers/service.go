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
	UpdateByEmail(subscriber entity.Subscribers) *newsletterError.ErrorCode
	Subscribe(subscriber entity.Subscribers) *newsletterError.ErrorCode
	Unsubscribe(subscriber entity.Subscribers) *newsletterError.ErrorCode
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

func (service *Service) UpdateByEmail(subscriber entity.Subscribers) *newsletterError.ErrorCode {
	err := service.Repo.UpdateByEmail(subscriber)

	if err != nil {
		return convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
	}

	return nil
}

func (service *Service) Subscribe(subscriber entity.Subscribers) *newsletterError.ErrorCode {

	resSubscribe, err := service.UseCase.FindByEmail(subscriber.Email)
	if err != nil {
		return convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
	}

	if len(resSubscribe) == 0 {
		newSubscribe := entity.Subscribers{
			Email: subscriber.Email,
			Name:  subscriber.Name,
		}

		errInsert := service.UseCase.Insert(newSubscribe)
		if errInsert != nil {
			return convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
		}

	} else {
		subscriber.IsSubscribed = true
		errInsert := service.UseCase.UpdateByEmail(subscriber)
		if errInsert != nil {
			return convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
		}
	}

	return nil
}

func (service *Service) Unsubscribe(subscriber entity.Subscribers) *newsletterError.ErrorCode {

	resSubscribe, err := service.UseCase.FindByEmail(subscriber.Email)
	if err != nil {
		return convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
	}

	if len(resSubscribe) == 0 {
		return convert.ValueToErrorCodePointer(newsletterError.DataNotFound)
	} else {
		newSubscribe := entity.Subscribers{
			Email:        subscriber.Email,
			Name:         subscriber.Name,
			IsSubscribed: false,
		}
		errInsert := service.UseCase.UpdateByEmail(newSubscribe)
		if errInsert != nil {
			return convert.ValueToErrorCodePointer(newsletterError.InternalServerError)
		}
	}

	return nil
}
