package main

import (
	"context"
	"database/sql"
	"newsletter/src/cmd/config"
	"newsletter/src/pkg/utils/logger"
	router "newsletter/src/routers"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "newsletter/src/routers/docs"

	_ "github.com/denisenkom/go-mssqldb"
)

const (
	defaultPort   = "8000"
	defaultAppEnv = "LOCAL"
)

func main() {
	config := config.New()
	fmt.Println("Stage", config.Stage)

	var (
		appPort  = getEnvString("PORT", defaultPort)
		httpAddr = flag.String("http.addr", ":"+appPort, "HTTP listen address")
	)
	log.Println("PORT: " + appPort)

	dbPOrt, _ := strconv.Atoi(config.DBPort)
	dbConnection, _ := connectDatabase(DBConnectURL{
		UserName: config.DBUserName,
		Password: config.DBPassword,
		DBHost:   config.DBHost,
		Port:     dbPOrt,
		DbName:   config.DBName,
	})

	logs, errLog := logger.NewELK(logger.ELKConnect{
		URL:      config.ELSURL,
		UserName: config.ELSUsername,
		Password: config.ELSPassword,
		Index:    config.ELSIndex,
		Stage:    config.Stage,
	})

	if errLog != nil {
		fmt.Println("Connect ELS Fail:", errLog)
	}
	if logs != nil {
		fmt.Println("Connect ELS Success", *logs)
	}

	routerConfig := router.RouterConfig{
		DB:     dbConnection,
		Logs:   logs,
		Config: config,
	}

	routers := router.InitRouter(routerConfig)

	server := httpServer(*httpAddr, routers)
	startServer(server)

	waitingForSignal(os.Interrupt, syscall.SIGTERM)

	fmt.Println("The service is shutting down...")

	forceShutdownAfter(server, time.Second*30)

	fmt.Println("terminated...")

	os.Exit(0)

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
