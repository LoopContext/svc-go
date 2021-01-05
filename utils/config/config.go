package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/loopcontext/msgcat"
	"github.com/loopcontext/svc-go/database"
)

var (
	defaultResponseContentType = "application/json"
	ErrDBMissingDSN            = fmt.Errorf("Missing DSN")
	ErrDBMissingPORT           = fmt.Errorf("Missing PORT")
)

type Config struct {
	Debug          bool
	Server         Server
	MessageCatalog msgcat.Config
	DB             database.DBConfig
}

type Server struct {
	Host                string
	Port                int
	ResponseContentType string
}

func (srvcfg *Server) BuildServerAddr() string {
	return srvcfg.Host + ":" + strconv.Itoa(srvcfg.Port)
}

func InitConfig() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, ErrDBMissingPORT
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, ErrDBMissingDSN
	}
	responseContentType := os.Getenv("RESPONSE_CONTENT_TYPE")
	if responseContentType == "" {
		responseContentType = defaultResponseContentType
	}
	return &Config{
		Debug: os.Getenv("DEBUG") == "TRUE",
		Server: Server{
			Host:                os.Getenv("HOST"),
			Port:                port,
			ResponseContentType: responseContentType,
		},
		MessageCatalog: msgcat.Config{
			CtxLanguageKey: os.Getenv("MESSAGE_CATALOG_LANGUAGE_KEY"),
			ResourcePath:   os.Getenv("MESSAGE_CATALOG_RESOURCE_PATH"),
		},
		DB: database.DBConfig{
			DSN:         dsn,
			Automigrate: os.Getenv("DB_AUTOMIGRATE") == "TRUE",
			Debug:       os.Getenv("DEBUG") == "TRUE",
		},
	}, nil
}
