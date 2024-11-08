package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	customtypes "winqroo/pkg/customTypes"
	"winqroo/pkg/utils"
)

func (h *UserAuthenticationHandler) UserLogoutHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	// userSession, ok := r.Context().Value(middlewares.UserSessionKey).(customtypes.UserClaim)
	// if !ok {
	// 	utils.SendHandlerCustomErrResponse(
	// 		w,
	// 		utils.NewCustomError(
	// 			fmt.Errorf("error in fetching user session"),
	// 			utils.ErrCodes.Common.ErrCodeInternalServerError,
	// 		),
	// 		http.StatusInternalServerError,
	// 	)
	// 	return
	// }

	// eventID := chi.URLParam(r, "eventID")
	// if eventID == "" {
	// 	utils.SendHandlerCustomErrResponse(
	// 		w,
	// 		utils.NewCustomError(
	// 			fmt.Errorf("eventID is a required param"),
	// 			utils.ErrCodes.Common.ErrCodeInternalServerError,
	// 		),
	// 		http.StatusBadRequest,
	// 	)
	// 	return
	// }

	// var putRewardableEventRequestModel loyaltyCustomTypes.PutRewardableEventRequestModel
	// if err := json.NewDecoder(r.Body).Decode(&putRewardableEventRequestModel); err != nil {
	// 	utils.SendHandlerCustomErrResponse(
	// 		w,
	// 		utils.NewCustomError(
	// 			fmt.Errorf("failed to json decode: %w", err),
	// 			utils.ErrCodes.Shared.ErrCodeInternalServerError,
	// 		),
	// 		http.StatusInternalServerError,
	// 	)
	// 	return
	// }

	// customErr := h.Service(
	// 	r.Context(),
	// 	&userSession,
	// )
	// if customErr != nil {
	// 	status := map[utils.CustomErrorCode]int{
	// 		utils.ErrCodes.Loyalty.RewardableEventNotFound: http.StatusNotFound,
	// 		utils.ErrCodes.Shared.ErrRestrictedAccess:      http.StatusUnauthorized,
	// 	}
	// 	if code, found := status[customErr.ErrCode]; found {
	// 		utils.SendHandlerCustomErrResponse(w, customErr, code)
	// 		return
	// 	}
	// 	utils.SendHandlerCustomErrResponse(w, customErr, http.StatusInternalServerError)
	// 	return
	// }
	cookie := http.Cookie{
		Name:     "user-jwt-token",
		Value:    "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusAccepted)
	// w.Header().Set(utils.ContentTypeHeader, utils.ApplicationJSON)
	// if err := json.NewEncoder(w).Encode(rewardableEvent); err != nil {
	// 	utils.SendHandlerCustomErrResponse(
	// 		w,
	// 		utils.NewCustomError(
	// 			fmt.Errorf("failed to json marshal: %w", err),
	// 			utils.ErrCodes.Common.ErrCodeInternalServerError,
	// 		),
	// 		http.StatusInternalServerError,
	// 	)
	// 	return
	// }
}

func (h *UserAuthenticationHandler) UserSignupHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var userSignup customtypes.UserSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&userSignup); err != nil {
		utils.SendHandlerCustomErrResponse(
			w,
			utils.NewCustomError(
				fmt.Errorf("failed to json decode: %w", err),
				utils.ErrCodes.Common.ErrCodeInternalServerError,
			),
			http.StatusInternalServerError,
		)
		return
	}

	customErr := h.Service.RegisterNewUser(r.Context(),userSignup)
	if customErr != nil {
		utils.SendHandlerCustomErrResponse(w, customErr, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserAuthenticationHandler) UserLoginHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var userLogin customtypes.UserLoginRequestModel
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		utils.SendHandlerCustomErrResponse(
			w,
			utils.NewCustomError(
				fmt.Errorf("failed to json decode: %w", err),
				utils.ErrCodes.Common.ErrCodeInternalServerError,
			),
			http.StatusInternalServerError,
		)
		return
	}

	userJwtToken, customErr := h.Service.LoginUser(r.Context(), userLogin.UserType, userLogin.UserEmailId, userLogin.Password)
	if customErr != nil {
		status := map[utils.CustomErrorCode]int{
			utils.ErrCodes.Common.ErrCodeRecordNotFound: http.StatusNotFound,
			utils.ErrCodes.Common.ErrUnAuthorised:       http.StatusUnauthorized,
		}
		if code, found := status[customErr.ErrCode]; found {
			utils.SendHandlerCustomErrResponse(w, customErr, code)
			return
		}
		utils.SendHandlerCustomErrResponse(w, customErr, http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "user-jwt-token",
		Value:    userJwtToken,
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusAccepted)
}
