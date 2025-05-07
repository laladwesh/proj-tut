package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status_error"`
	Error  string
}

const (
	StatusOk    = "success"
	StatusError = "error"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(err validator.ValidationErrors) Response {
	var errMsgs []string
	for _, fieldErr := range err {
		switch fieldErr.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fieldErr.Field()+" is required")
		default:
			errMsgs = append(errMsgs, fieldErr.Field()+" is not valid")

			// case "min":
			// 	errMsgs = append(errMsgs, fieldErr.Field()+" must be at least "+fieldErr.Param())
			// case "max":
			// 	errMsgs = append(errMsgs, fieldErr.Field()+" must be at most "+fieldErr.Param())
			// case "len":
			// 	errMsgs = append(errMsgs, fieldErr.Field()+" must be exactly "+fieldErr.Param())
			// case "oneof":
			// 	errMsgs = append(errMsgs, fieldErr.Field()+" must be one of "+fieldErr.Param())
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
