package main

import (
	"crypto/tls"
	"database/sql"
	"net/url"
	configs "subscribetool/src/cmd/config"
	"subscribetool/src/pkg/subscribers"
	"subscribetool/src/pkg/utils/email"
	"subscribetool/src/pkg/utils/logger"

	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/gomail.v2"
)

func main() {
	config := configs.GetConfig()

	dbConnection, _ := connectDatabase(DBConnectURL{
		UserName: config.MSSQL.MssqlUsername,
		Password: config.MSSQL.MssqlPassword,
		DBHost:   config.MSSQL.MssqlHost,
		Port:     config.MSSQL.MssqlPort,
		DbName:   config.MSSQL.MSsqlName,
	})

	var emailServerTLSSkipVerify bool
	if config.EmailServer.EmailSMTPTLSSkipVarify == "true" {
		emailServerTLSSkipVerify = true
	} else {
		emailServerTLSSkipVerify = false
	}

	dialerMailServer, _ := makeDialerMailServer(EmailSMTPConnect{
		EmailSMTPHost:          config.EmailServer.EmailSMTPHost,
		EmailSmtpPort:          config.EmailServer.EmailSMTPPort,
		EmailSMTPEmailSender:   config.EmailServer.EmailSMTPMailSender,
		EmailSMTPPassword:      config.EmailServer.EmailSMTPPassword,
		EmailSMTPTLSSkipVerify: emailServerTLSSkipVerify,
	})

	logs, errLog := logger.NewELK(logger.ELKConnect{
		URL:      config.Elastic.Host,
		UserName: config.Elastic.UserName,
		Password: config.Elastic.Password,
		Index:    config.Elastic.Index,
		Stage:    "dev",
	})

	if errLog != nil {
		fmt.Println("Connect ELS Fail:", errLog)
	}
	if logs != nil {
		fmt.Println("Connect ELS Success", *logs)
	}

	//repository
	subscribersRepository := subscribers.NewRepository("TB_TRN_Subscribers", dbConnection, logs)

	// Service
	email.FromEMailSender = config.EmailServer.EmailSMTPMailSender
	utilsEmailService := email.NewService(dialerMailServer)
	serviceParam := subscribers.ServiceParam{
		UtilsEmailService: utilsEmailService,
		Repo:              subscribersRepository,
		Logs:              logs,
	}

	subscribersService := subscribers.NewService(serviceParam)

	subscribersService.SentEmail()

	fmt.Println("End of Process")
}

type DBConnectURL struct {
	UserName string
	Password string
	DBHost   string
	Port     int
	DbName   string
}

func connectDatabase(connectionInfo DBConnectURL) (*sql.DB, error) {
	query := url.Values{}

	query.Add("database", connectionInfo.DbName)
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(connectionInfo.UserName, connectionInfo.Password),
		Host:     fmt.Sprintf("%s:%d", connectionInfo.DBHost, connectionInfo.Port),
		RawQuery: query.Encode(),
	}
	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		fmt.Println("Sql fail !!")
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Sql success !!")
	return db, nil
}

type EmailSMTPConnect struct {
	EmailSMTPHost          string
	EmailSmtpPort          int
	EmailSMTPEmailSender   string
	EmailSMTPPassword      string
	EmailSMTPTLSSkipVerify bool
}

func makeDialerMailServer(connectionMailServer EmailSMTPConnect) (*gomail.Dialer, error) {

	mailer := gomail.NewDialer(
		connectionMailServer.EmailSMTPHost,
		connectionMailServer.EmailSmtpPort,
		connectionMailServer.EmailSMTPEmailSender,
		connectionMailServer.EmailSMTPPassword,
	)
	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: connectionMailServer.EmailSMTPTLSSkipVerify}

	fmt.Println("Make Dialer For Mail Server Success !!")

	return mailer, nil
}
