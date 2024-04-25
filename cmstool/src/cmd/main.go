package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"strconv"
	configs "subscribetool/src/cmd/config"
	"subscribetool/src/pkg/subscribers"
	"subscribetool/src/pkg/utils/email"
	"subscribetool/src/pkg/utils/logger"

	// "subscribetool/src/cmd/config"
	// "subscribetool/src/pkg/utils/logger"
	// router "subscribetool/src/routers"
	// "flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	// "strconv"
	// "syscall"
	"time"

	// _ "subscribetool/src/routers/docs"

	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/gomail.v2"
)

const (
	defaultPort   = "8000"
	defaultAppEnv = "LOCAL"
)

func main() {
	config := configs.GetConfig()
	fmt.Println("config : ", config)

	dbPOrt, _ := strconv.Atoi(config.MSSQL.MssqlHost)
	dbConnection, _ := connectDatabase(DBConnectURL{
		UserName: config.MSSQL.MssqlUsername,
		Password: config.MSSQL.MssqlPassword,
		DBHost:   config.MSSQL.MssqlHost,
		Port:     dbPOrt,
		DbName:   config.MSSQL.MSsqlName,
	})

	emailServerPort, _ := strconv.Atoi(config.EmailServer.EmailSMTPTLSSkipVarify)

	var emailServerTLSSkipVerify bool
	if config.EmailServer.EmailSMTPTLSSkipVarify == "true" {
		emailServerTLSSkipVerify = true
	} else {
		emailServerTLSSkipVerify = false
	}

	dialerMailServer, _ := makeDialerMailServer(EmailSMTPConnect{
		EmailSMTPHost:          config.EmailServer.EmailSMTPHost,
		EmailSmtpPort:          emailServerPort,
		EmailSMTPEmailSender:   config.EmailServer.EmailSMTPMailSender,
		EmailSMTPPassword:      config.EmailServer.EmailSMTPPassword,
		EmailSMTPTLSSkipVerify: emailServerTLSSkipVerify,
	})

	logs, errLog := logger.NewELK(logger.ELKConnect{
		URL:      config.Elastic.Host,
		UserName: config.Elastic.UserName,
		Password: config.Elastic.Password,
		Index:    config.Elastic.Index,
		Stage:    "local",
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

func getEnvString(env, fallback string) string {
	result := os.Getenv(env)
	if result == "" {
		return fallback
	}
	return result
}

func httpServer(httpAddr string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}
}

func startServer(server *http.Server) {
	ch := make(chan error, 1)

	go func() {
		ch <- server.ListenAndServe()
	}()

	select {
	case err := <-ch:
		log.Fatal(err)
	default:
		log.Println("The service is ready to listen and serve.")
	}
}

func forceShutdownAfter(server *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	server.Shutdown(ctx)
}

func waitingForSignal(sig ...os.Signal) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, sig...)

	s := <-stop
	log.Println("Got signal ", s.String())
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
