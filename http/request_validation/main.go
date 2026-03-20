package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	helperfuncs "github.com/baizhigit/go-snippets/http/helper_funcs"
	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validate: v}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

var validate = validator.New()

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Age      int    `json:"age" validate:"gte=18,lte=120"`
}

func (r *CreateUserRequest) Validate() error {
	return validate.Struct(r)
}

func handleCreateUser(v *Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			helperfuncs.RespondJSON(w, http.StatusBadRequest, APIError{
				Code:    "INVALID_JSON",
				Message: "Invalid request body",
			})
			return
		}

		if err := v.Validate(&req); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {

				fields := make(map[string]string)
				for _, e := range ve {
					fields[e.Field()] = getErrorMessage(e)
				}

				helperfuncs.RespondJSON(w, http.StatusBadRequest, APIError{
					Code:    "VALIDATION_ERROR",
					Message: "Validation failed",
					Fields:  fields,
				})
				return
			}
		}

		// success
		helperfuncs.RespondJSON(w, http.StatusCreated, map[string]string{
			"status": "ok",
		})
	}
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be >= %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be <= %s", fe.Param())
	default:
		return "Invalid value"
	}
}
