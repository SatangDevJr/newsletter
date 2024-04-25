package configs

import (
	//"encoding/asn1"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Env         string
	MSSQL       MSSQL
	Elastic     Elastic
	EmailServer EmailServer
}

type MSSQL struct {
	MssqlUsername string
	MssqlPassword string
	MssqlPort     int
	MssqlHost     string
	MSsqlName     string
}

type Elastic struct {
	Host     string
	Index    string
	UserName string
	Password string
}

type EmailServer struct {
	EmailSMTPHost          string
	EmailSMTPPort          int
	EmailSMTPMailSender    string
	EmailSMTPPassword      string
	EmailSMTPTLSSkipVarify string
}

func GetConfig(params ...string) Configuration {
	configuration := Configuration{}
	env := os.Getenv("env")
	if env == "" {
		env = "local"
	}

	if len(params) > 0 {
		env = params[0]
	}

	ex, err := os.Executable()
	if err != nil {
		fmt.Printf("error %+v\n", err)
	}

	exPath := filepath.Dir(ex)
	fileName := fmt.Sprintf("%s/configs/%s_config.json", exPath, env)

	err = gonfig.GetConf(fileName, &configuration)
	if err != nil {
		fileName = fmt.Sprintf("./configs/%s_config.json", env)
		err = gonfig.GetConf(fileName, &configuration)
		if err != nil {
			fmt.Println(err)
		}
	}
	return configuration
}
