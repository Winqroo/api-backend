package routes

import (
	"net/http"

	"winqroo/middlewares"
	apiRoutes "winqroo/routes/api"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/go-chi/chi/v5"
	// "go.uber.org/zap"
)

func NewRoutes(ddb *dynamodb.Client, ses *ses.Client) http.Handler {
	router := chi.NewRouter()

	// Middlewares
	authMiddleware := middlewares.NewAuthMiddleware()
	// TODO: logger for next version
	// loggerMiddleWare := middlewares.NewRequestLoggerMiddleware().RequestLogMiddleware(logger)
	// router.Use(loggerMiddleWare)

	// Base route for health check
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Base route for our use
	router.Route("/api/v1", loadApiRoutes(authMiddleware, ddb, ses))

	return router
}

func loadApiRoutes(
	authMiddleware *middlewares.AuthMiddleware,
	ddb *dynamodb.Client,
	ses *ses.Client,
) func(router chi.Router) {
	return func(r chi.Router) {
		r.Route("/user", apiRoutes.UserRoutes(authMiddleware, ddb, ses))
		// r.Route("/tasks", apiRoutes.TaskRoutes(authMiddleware, ddb))
		// r.Route("/creator", apiRoutes.CreatorRoutes(authMiddleware, ddb))
		// r.Route("/brand", apiRoutes.BrandRoutes(authMiddleware, ddb))
	}
}
