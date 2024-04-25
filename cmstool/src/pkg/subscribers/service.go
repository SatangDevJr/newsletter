package subscribers

import (
	"fmt"
	"subscribetool/src/pkg/entity"
	"subscribetool/src/pkg/utils/convert"
	"subscribetool/src/pkg/utils/email"
	subscribetoolError "subscribetool/src/pkg/utils/error"
	"subscribetool/src/pkg/utils/logger"
)

type UseCase interface {
	GetAllSubscribers() ([]entity.Subscribers, *subscribetoolError.ErrorCode)
	SentEmail() *subscribetoolError.ErrorCode
}

type Service struct {
	UseCase
	UtilsEmailService email.UseCase
	Repo              Repository
	Logs              logger.Logger
}

func NewService(serviceParam ServiceParam) *Service {
	service := &Service{
		UtilsEmailService: serviceParam.UtilsEmailService,
		Repo:              serviceParam.Repo,
		Logs:              serviceParam.Logs,
	}
	service.UseCase = service
	return service
}

func (service *Service) GetAllSubscribers() ([]entity.Subscribers, *subscribetoolError.ErrorCode) {
	res, err := service.Repo.GetAllSubscribers()

	if err != nil {
		return nil, convert.ValueToErrorCodePointer(subscribetoolError.InternalServerError)
	}

	if res == nil {
		return nil, convert.ValueToErrorCodePointer(subscribetoolError.DataNotFound)
	}

	return res, nil
}

func (service *Service) SentEmail() *subscribetoolError.ErrorCode {

	resGetAllSubscribers, errGetAllSubscribers := service.UseCase.GetAllSubscribers()
	if errGetAllSubscribers != nil {
		return convert.ValueToErrorCodePointer(subscribetoolError.InternalServerError)
	}

	if len(resGetAllSubscribers) == 0 {
		return convert.ValueToErrorCodePointer(subscribetoolError.DataNotFound)
	}

	for _, value := range resGetAllSubscribers {

		emailTarget := []string{value.Email}
		fmt.Println("emailTarget : ", emailTarget)
		mailInfo := email.SentMailContent{
			To:      emailTarget,
			Supject: "Test sent mail for subscribers",
			Body:    "This email sent for notification. Test sent mail for subscribers",
		}

		errUtilsEmail := service.UtilsEmailService.Send(mailInfo)
		if errUtilsEmail != nil {
			return convert.ValueToErrorCodePointer(subscribetoolError.InternalServerError)
		}

	}

	fmt.Println("sent mail to subscribers success")

	return nil
}
