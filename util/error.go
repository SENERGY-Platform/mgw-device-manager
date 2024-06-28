package util

import (
	"errors"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"net/http"
)

func GetStatusCode(err error) int {
	var nfe *lib_model.NotFoundError
	if errors.As(err, &nfe) {
		return http.StatusNotFound
	}
	var iie *lib_model.InvalidInputError
	if errors.As(err, &iie) {
		return http.StatusBadRequest
	}
	var rbe *lib_model.ResourceBusyError
	if errors.As(err, &rbe) {
		return http.StatusConflict
	}
	var ie *lib_model.InternalError
	if errors.As(err, &ie) {
		return http.StatusInternalServerError
	}
	return 0
}
