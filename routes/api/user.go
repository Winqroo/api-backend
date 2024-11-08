package api

import (
	"winqroo/middlewares"
	userHandler "winqroo/pkg/handlers/users"
	userRepos "winqroo/pkg/repositories/users"
	userServices "winqroo/pkg/services/users"

	// "winqroo/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(
	authMiddleware *middlewares.AuthMiddleware,
	ddb *dynamodb.Client,
	ses *ses.Client,
) func(router chi.Router) {
	return func(r chi.Router) {
		userProfileRepo := userRepos.NewUserProfileRepo(ddb)
		emailVerificationRepo := userRepos.NewUserEmailVerificationRepo(ddb, ses)

		userAuthenticationService := userServices.NewUserAuthenticationService(userProfileRepo)
		userVerificationService := userServices.NewUserVerificationService(emailVerificationRepo)
		userInteractionsService := userServices.NewUserInteractionsService(userProfileRepo)

		userAuthenticationHandler := userHandler.NewUserAuthenticationHandler(userAuthenticationService)
		userVerificationHandler := userHandler.NewUserVerificationHandler(userVerificationService)
		userInteractionsHandler := userHandler.NewUserInteractionsHandler(userInteractionsService)

		r.Route("/auth", func(sr chi.Router) {
			sr.Route("/otp", func(subRouter chi.Router) {
				subRouter.Post(
					"/register",
					userVerificationHandler.GetOtpToRegisterHandler,
				)
				subRouter.Post(
					"/resend",
					userVerificationHandler.ResendOtpHandler,
				)
				subRouter.Post(
					"/verify",
					userVerificationHandler.VerifyOtpHandler,
				)
			})
			sr.Post("/login", userAuthenticationHandler.UserLoginHandler)
			sr.Put("/signup", userAuthenticationHandler.UserSignupHandler)
			sr.With(authMiddleware.UserSessionJwtAuthMiddleware).Post("/logout", userAuthenticationHandler.UserLogoutHandler)
		})

		r.Route("/service", func(sr chi.Router) {
			sr.Use(authMiddleware.UserSessionJwtAuthMiddleware)
			sr.Get("/{userId}", userInteractionsHandler.GetUserByIDHandler)
		})
	}
}
