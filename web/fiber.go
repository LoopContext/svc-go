package server

import (
	"net/http"

	"github.com/loopcontext/svc-go/web/handlers/clients"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/loopcontext/msgcat"
	"github.com/loopcontext/svc-go/database"
	"github.com/loopcontext/svc-go/database/repo"
	"github.com/loopcontext/svc-go/services/managers"
	"github.com/loopcontext/svc-go/utils"
	"github.com/loopcontext/svc-go/utils/config"
	"github.com/loopcontext/svc-go/utils/logger"
	"github.com/loopcontext/svc-go/web/handlers"
	"github.com/loopcontext/svc-go/web/middleware"
)

func InitFiberServer(cfg *config.Config, log *logger.Logger,
	catalog *msgcat.MessageCatalog, db *database.DB) (app *fiber.App, err error) {
	app = fiber.New()
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} - ${ip} - ${pid} - ${locals:requestid} - ${ua} - ${status} ${method} ${path}​ [url: ${url}]​​ \n",
	}))

	err = setupRoutes(app, cfg, log, *catalog, db)

	return app, err
}

func setupRoutes(app *fiber.App, cfg *config.Config, log *logger.Logger, catalog msgcat.MessageCatalog, db *database.DB) (err error) {
	// Default handler
	pingHandler := adaptor.HTTPHandlerFunc(handlers.Ping)
	// Repos
	authRepo := repo.NewAuthRepoSvc(db, *log, catalog)
	clientRepo := repo.NewClientRepoSvc(db)
	// HTTP server's request and response service
	requestHelperSvc := helpers.NewRequestHelperSvc(*log, catalog)
	responseHelperHelperSvc := helpers.NewResponseHelperSvc(*log, catalog, cfg.Server.ResponseContentType)
	// HTTP server's helpers
	authHelperSvc := helpers.NewAuthHelperSvc(catalog)
	// Auth middleware
	authMiddleware := middleware.NewAuthMiddlewareSvc(*log, catalog, authRepo, responseHelperHelperSvc)
	apiKeyMiddleware := HTTPMiddleware(authMiddleware.CheckUserAPIKey)
	// Clients
	clientHandler := clients.NewClientHandlerSvc(*log, catalog, clientRepo, authHelperSvc, requestHelperSvc, responseHelperHelperSvc)
	clientRead := adaptor.HTTPHandlerFunc(clientHandler.Read)
	clientList := adaptor.HTTPHandlerFunc(clientHandler.List)

	app.Get("/", pingHandler)
	app.Get("/ping", pingHandler)
	app.Get("/health", pingHandler)

	api := app.Group("/api", apiKeyMiddleware)

	// Clients
	clientRoute := api.Group("/clients")
	clientRoute.Get("/", clientList)
	clientRoute.Get("/:id", clientRead)

	return err
}

// HTTPMiddleware wraps net/http middleware to fiber middleware
func HTTPMiddleware(mw func(http.Handler) http.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var next bool
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next = true
			// Convert again in case request may modify by middleware
			c.Locals(string(utils.CurrrentUserCtxKey), r.Context().Value(utils.CurrrentUserCtxKey))
			c.Request().Header.SetMethod(r.Method)
			c.Request().SetRequestURI(r.RequestURI)
			c.Request().SetHost(r.Host)
			for key, val := range r.Header {
				for _, v := range val {
					c.Request().Header.Set(key, v)
				}
			}
		})
		_ = adaptor.HTTPHandler(mw(nextHandler))(c)
		if next {
			return c.Next()
		}
		return nil
	}
}
