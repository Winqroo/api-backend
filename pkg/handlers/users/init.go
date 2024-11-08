package users

import userServices "winqroo/pkg/services/users"

type (
	UserAuthenticationHandler struct {
		Service *userServices.UserAuthenticationService
	}
	UserVerificationHandler struct {
		Service *userServices.UserVerificationService
	}
	UserInteractionsHandler struct {
		Service *userServices.UserInteractionsService
	}
)

func NewUserAuthenticationHandler(
	s *userServices.UserAuthenticationService,
) *UserAuthenticationHandler {
	return &UserAuthenticationHandler{
		Service: s,
	}
}

func NewUserVerificationHandler(
	s *userServices.UserVerificationService,
) *UserVerificationHandler {
	return &UserVerificationHandler{
		Service: s,
	}
}

func NewUserInteractionsHandler(
	s *userServices.UserInteractionsService,
) *UserInteractionsHandler {
	return &UserInteractionsHandler{
		Service: s,
	}
}
