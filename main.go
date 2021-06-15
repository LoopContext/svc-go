package main

import (
	"github.com/loopcontext/msgcat"
	"github.com/loopcontext/svc-go/database"
	"github.com/loopcontext/svc-go/utils/config"
	"github.com/loopcontext/svc-go/utils/logger"
	server "github.com/loopcontext/svc-go/web"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg.Debug)
	catalog, err := msgcat.NewMessageCatalog(cfg.MessageCatalog)
	if err != nil {
		log.SendFatal(err)
	}
	db, err := database.NewDB(cfg.DB, log)
	if err != nil {
		log.SendFatal(err)
	}
	appsrv, err := server.InitFiberServer(cfg, log, &catalog, db)
	if err != nil {
		log.SendFatal(err)
	}
	log.SendFatal(appsrv.Listen(cfg.Server.BuildServerAddr()))
}
