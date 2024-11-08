package middlewares

import (
	"context"
	"fmt"
	"net/http"
	customtypes "winqroo/pkg/customTypes"
	"winqroo/pkg/utils"
)

const (
	UserSessionKey = "UserSession"
)

type AuthMiddleware struct {
	jwtAuthSystem *utils.JwtAuthSystemUtils
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		jwtAuthSystem: utils.NewJwtAuthSystemUtils(),
	}
}

func (a *AuthMiddleware) UserSessionJwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In ApiAuthMiddleware")
		jwtToken, err := r.Cookie("user-jwt-token")
		if err != nil {
			utils.SendHandlerCustomErrResponse(
				w,
				utils.NewCustomError(
					fmt.Errorf("missing required cookie"),
					utils.ErrCodes.Common.ErrCodeBadRequest,
				),
				http.StatusBadRequest,
			)

			return
		}

		var userSession *customtypes.UserClaim 
		err = a.jwtAuthSystem.ValidateSessionJWT(jwtToken.String(), userSession)
		if err != nil {
			utils.SendHandlerCustomErrResponse(
				w,
				utils.NewCustomError(
					fmt.Errorf("authorization failed: %v", err),
					utils.ErrCodes.Common.ErrUnAuthorised,
				),
				http.StatusUnauthorized,
			)
			return
		}

		ctx := context.WithValue(r.Context(), UserSessionKey, userSession)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
