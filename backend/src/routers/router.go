package router

import (
	"database/sql"
	"fmt"
	"net/http"
	"newsletter/src/api/middleware"
	"os"

	"newsletter/src/cmd/config"

	subscribers "newsletter/src/pkg/subscribers"
	"newsletter/src/pkg/utils/logger"

	subscribersHandler "newsletter/src/api/subscribers/handler"

	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	defaultPort   = "8000"
	defaultAppEnv = "LOCAL"
)

type RouterConfig struct {
	DB     *sql.DB
	Logs   *logger.ELK
	Config config.Configuration
}

func InitRouter(routerConfig RouterConfig) http.Handler {
	fmt.Println("InitRouter :", routerConfig)

	/* Repository */
	subscribersRepository := subscribers.NewRepository("TB_TRN_Subscribers", routerConfig.DB, routerConfig.Logs)

	/* Service */
	subscribersService := subscribers.NewService(subscribersRepository, routerConfig.Logs)

	/* Handler */
	subscribersHandlerParam := subscribersHandler.HandlerParam{
		Service: subscribersService,
		Logs:    routerConfig.Logs,
	}
	subscribersHandler := subscribersHandler.MakeSubscribersHandler(subscribersHandlerParam)

	/* Router */
	middleware := middleware.NewMiddleware(routerConfig.Logs)
	router := mux.NewRouter()
	router.Use(middleware.Recover)
	router.HandleFunc("/version", versionHandler)
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	subscribers := router.PathPrefix("/subscribers").Subrouter()
	subscribers.HandleFunc("", http.HandlerFunc(subscribersHandler.GetAllSubscribers)).Methods("GET")

	return router
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	appVersion := getEnvString("APP_VERSION", defaultPort)
	fmt.Fprintln(w, appVersion)
}

func getEnvString(env, fallback string) string {
	result := os.Getenv(env)
	if result == "" {
		return fallback
	}
	return result
}

type DBConnectURL struct {
	UserName string
	Password string
	DBHost   string
	Port     int
	DbName   string
}
