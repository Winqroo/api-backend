package utils

import (
	"encoding/json"
	"net/http"
)

const (
	ContentTypeHeader = "Content-Type"
	ApplicationJSON   = "application/json"
)

type (
	CustomErrorCode string

	commonErrCodes struct {
		ErrCodeBadRequest          CustomErrorCode
		ErrCodeInternalServerError CustomErrorCode
		ErrUnAuthorised            CustomErrorCode
		ErrRestrictedAccess        CustomErrorCode
		ErrCodeConflict            CustomErrorCode
		ErrCodeRecordNotFound      CustomErrorCode
	}

	customErrCodes struct {

	}

	ErrorCodes struct {
		Common commonErrCodes
		Custom customErrCodes
	}

	CustomError struct {
		Error   error           `json:"error"`
		ErrCode CustomErrorCode `json:"errCode"`
	}
)

var ErrCodes = &ErrorCodes{
	Common: commonErrCodes{
		ErrCodeBadRequest:          "BAD_REQUEST",
		ErrCodeInternalServerError: "INTERNAL_SERVER_ERROR",
		ErrUnAuthorised:            "UN_AUTHORIZED",
		ErrRestrictedAccess:        "ACCESS_DENIED",
		ErrCodeConflict:            "CONFLICT",
		ErrCodeRecordNotFound:      "RECORD_NOT_FOUND",
	},
}

func NewCustomError(err error, errCode CustomErrorCode) *CustomError {
	return &CustomError{
		Error:   err,
		ErrCode: errCode,
	}
}

func SendHandlerCustomErrResponse(w http.ResponseWriter, customErr *CustomError, status int) {
	customErrStr := struct {
		Error   string          `json:"error"`
		ErrCode CustomErrorCode `json:"errCode"`
	}{
		Error:   customErr.Error.Error(),
		ErrCode: customErr.ErrCode,
	}

	responseJSON, err := json.Marshal(customErrStr)
	if err != nil {
		w.Header().Set(ContentTypeHeader, ApplicationJSON)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}

	w.Header().Set(ContentTypeHeader, ApplicationJSON)
	w.WriteHeader(status)
	w.Write(responseJSON)
}

func (e *CustomError) Is(target *CustomError) bool {
	return e.Error == target.Error && e.ErrCode == target.ErrCode
}

func ReturnError(err error, errCode CustomErrorCode) *CustomError {
	// LogError(ctx, err)  //not for v1
	return NewCustomError(err, errCode)
}
