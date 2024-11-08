package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"winqroo/pkg/utils"
)

type otpData struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	AuthCode string `json:"authCode"`
}

func (h *UserVerificationHandler) GetOtpToRegisterHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var otp otpData
	if err := json.NewDecoder(r.Body).Decode(&otp); err != nil {
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

	data, err := h.Service.GetOtpToRegister(r.Context(), otp.Email)
	if err != nil {
		utils.SendHandlerCustomErrResponse(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(utils.ContentTypeHeader, utils.ApplicationJSON)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		utils.SendHandlerCustomErrResponse(
			w,
			utils.NewCustomError(
				fmt.Errorf("failed to json marshal: %w", err),
				utils.ErrCodes.Common.ErrCodeInternalServerError,
			),
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *UserVerificationHandler) ResendOtpHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var otp otpData
	if err := json.NewDecoder(r.Body).Decode(&otp); err != nil {
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

	data, err := h.Service.ResendOtp(r.Context(), otp.Email)
	if err != nil {
		utils.SendHandlerCustomErrResponse(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(utils.ContentTypeHeader, utils.ApplicationJSON)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		utils.SendHandlerCustomErrResponse(
			w,
			utils.NewCustomError(
				fmt.Errorf("failed to json marshal: %w", err),
				utils.ErrCodes.Common.ErrCodeInternalServerError,
			),
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *UserVerificationHandler) VerifyOtpHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var otp otpData
	if err := json.NewDecoder(r.Body).Decode(&otp); err != nil {
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

	if otp.Email == "" || otp.OTP == "" || otp.AuthCode == "" {
		utils.SendHandlerCustomErrResponse(
			w,
			utils.NewCustomError(
				fmt.Errorf("Incomplete data"),
				utils.ErrCodes.Common.ErrCodeBadRequest,
			),
			http.StatusBadRequest,
		)
		return
	}

	data, err := h.Service.VerifyOtp(r.Context(), otp.Email, otp.OTP, otp.AuthCode)
	if err != nil {
		utils.SendHandlerCustomErrResponse(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(utils.ContentTypeHeader, utils.ApplicationJSON)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		utils.SendHandlerCustomErrResponse(
			w,
			utils.NewCustomError(
				fmt.Errorf("failed to json marshal: %w", err),
				utils.ErrCodes.Common.ErrCodeInternalServerError,
			),
			http.StatusInternalServerError,
		)
		return
	}
}
