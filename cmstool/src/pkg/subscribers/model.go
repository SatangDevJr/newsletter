package subscribers

import (
	"subscribetool/src/pkg/utils/email"
	"subscribetool/src/pkg/utils/logger"
)

type ServiceParam struct {
	UtilsEmailService email.UseCase
	Repo              Repository
	Logs              logger.Logger
}
