package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/laladwesh/proj-tut/internal/storage"
	"github.com/laladwesh/proj-tut/internal/types"
	"github.com/laladwesh/proj-tut/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate the student struct (this is a placeholder, implement your own validation logic)
		validate := validator.New()
		if err := validate.Struct(student); err != nil {
			if verr, ok := err.(validator.ValidationErrors); ok {
				response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(verr))
			} else {
				response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			}
			return
		}

		lastid, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}
		slog.Info("user created successfully ", slog.String("userId", fmt.Sprint(lastid)))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastid})
		slog.Info("Creating a student", slog.String("method", r.Method), slog.String("path", r.URL.Path))
		response.WriteJSON(w, http.StatusCreated, map[string]string{"message": "student created"})
	}
}
