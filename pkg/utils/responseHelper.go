package utils

import (
	model "chat-app/pkg/models"
	"net/http"
)

func NewResponse(status int, data any, message string, success bool) model.ResponseModel {
	return model.ResponseModel{
		Status:  status,
		Data:    data,
		Message: message,
		Success: success,
	}
}

func StatusOK(data any) model.ResponseModel {
	return model.ResponseModel{
		Status:  http.StatusOK,
		Data:    data,
		Message: "OK",
		Success: true,
	}
}

func StatusInternalServerError() model.ResponseModel {
	return model.ResponseModel{
		Status:  http.StatusInternalServerError,
		Message: "Unhandled error occurred. Please try again later",
		Success: false,
	}
}

func StatusBadRequest(message string) model.ResponseModel {
	return model.ResponseModel{
		Status:  http.StatusBadRequest,
		Message: message,
		Success: false,
	}
}
