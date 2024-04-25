package email

import (
	"errors"
	"log"

	"gopkg.in/gomail.v2"
)

var (
	FromEMailSender string
)

type Service struct {
	Service          UseCase
	DialerMailServer *gomail.Dialer
}

type UseCase interface {
	Send(sentMailContent SentMailContent) error
}

type ServiceParam struct {
	DialerMailServer *gomail.Dialer
}

func NewService(dialerMailServer *gomail.Dialer) *Service {
	return &Service{
		DialerMailServer: dialerMailServer,
	}
}

func (service *Service) Send(sentMailContent SentMailContent) error {
	message := gomail.NewMessage()

	var emails []string
	if len(sentMailContent.To) != 0 && sentMailContent.To != nil {
		for _, mailTo := range sentMailContent.To {
			if mailTo != "" {
				emails = append(emails, mailTo)
			}
		}
	} else {
		errMailtoEmpty := errors.New("email to empty")
		return errMailtoEmpty
	}

	if len(emails) == 0 || emails == nil {
		errMailtoEmpty := errors.New("email to empty")
		return errMailtoEmpty
	}

	message.SetHeader("To", emails...)
	message.SetHeader("Subject", sentMailContent.Supject)
	message.SetBody("text/html", sentMailContent.Body)
	message.SetHeader("From", FromEMailSender)

	if err := service.DialerMailServer.DialAndSend(message); err != nil {
		log.Println("Email sent DialAndSend err : ", err)
		return err
	}

	log.Println("Email sent successfully!")
	return nil
}
