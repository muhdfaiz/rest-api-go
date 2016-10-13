package systems

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	validator "gopkg.in/go-playground/validator.v8"
)

type Validation struct {
	ErrorMessages map[string]string
	Validator     *validator.Validate
}

func (v *Validation) Validate(c *gin.Context, validationTags map[string]string) *ErrorData {
	v.initValidator()

	for validationKey, validationRules := range validationTags {

		validationRules := strings.Split(validationRules, ",")

		for _, validationRule := range validationRules {
			var message string

			validationRule := strings.Split(validationRule, "=")

			value := c.Query(validationKey)

			switch validationRule[0] {
			case "required":
				message = v.validateRequired(value, validationKey)
			case "uuid5":
				message = v.validateUUIDV5(value, validationKey)
			case "alpha":
				message = v.validateAlpha(value, validationKey)
			case "alphanum":
				message = v.validateAlphaNumeric(value, validationKey)
			case "numeric":
				message = v.validateNumeric(value, validationKey)
			case "min":
				message = v.validateMin(value, validationKey, validationRule[1])
			case "max":
				message = v.validateMax(value, validationKey, validationRule[1])
			case "email":
				message = v.validateEmail(value, validationKey)
			case "len":
				message = v.validateLength(value, validationKey, validationRule[1])
			case "time":
				message = v.validateTime(value, validationKey)
			}

			if message != "" {
				v.ErrorMessages[validationKey] = message
			}
		}
	}

	if len(v.ErrorMessages) > 0 {
		return &ErrorData{
			Error: &ErrorFormat{
				Status: strconv.Itoa(http.StatusUnprocessableEntity),
				Code:   QueryStringValidationFailed,
				Title:  TitleValidationError,
				Detail: v.ErrorMessages,
			},
		}
	}

	return nil

}

func (v *Validation) initValidator() {
	config := &validator.Config{TagName: "validate"}
	v.Validator = validator.New(config)
	v.ErrorMessages = make(map[string]string)
}

func (v *Validation) validateRequired(value string, validationKey string) string {
	err := v.Validator.Field(value, "required")

	if err != nil {
		return fmt.Sprintf(ErrorValidationRequired, validationKey)
	}

	return ""
}

func (v *Validation) validateUUIDV5(value string, validationKey string) string {
	if value != "" {
		err := v.Validator.Field(value, "uuid5")

		if err != nil {
			return fmt.Sprintf(ErrorValidationUUID, validationKey)
		}
	}

	return ""
}

func (v *Validation) validateAlpha(value string, validationKey string) string {
	if value != "" {
		err := v.Validator.Field(value, "alpha")

		if err != nil {
			return fmt.Sprintf(ErrorValidationAlpha, validationKey)
		}
	}

	return ""
}

func (v *Validation) validateAlphaNumeric(value string, validationKey string) string {
	if value != "" {
		err := v.Validator.Field(value, "alphanum")

		if err != nil {
			return fmt.Sprintf(ErrorValidationAlphaNum, validationKey)
		}
	}

	return ""
}

func (v *Validation) validateNumeric(value string, validationKey string) string {
	if value != "" {
		err := v.Validator.Field(value, "numeric")

		if err != nil {
			return fmt.Sprintf(ErrorValidationNumeric, validationKey)
		}
	}

	return ""
}

func (v *Validation) validateMin(value string, validationKey string, minLength string) string {
	if value != "" {
		err := v.Validator.Field(value, "min"+"="+minLength)

		if err != nil {
			return fmt.Sprintf(ErrorValidationMin, validationKey, minLength)
		}
	}

	return ""
}

func (v *Validation) validateMax(value string, validationKey string, maxLength string) string {
	if value != "" {
		err := v.Validator.Field(value, "max"+"="+maxLength)

		if err != nil {
			return fmt.Sprintf(ErrorValidationMax, validationKey, maxLength)
		}
	}

	return ""
}

func (v *Validation) validateEmail(value string, validationKey string) string {
	if value != "" {
		err := v.Validator.Field(value, "email")

		if err != nil {
			return fmt.Sprintf(ErrorValidationEmail, validationKey)
		}
	}

	return ""
}

func (v *Validation) validateLength(value string, validationKey string, length string) string {
	if value != "" {
		err := v.Validator.Field(value, "len"+"="+length)

		if err != nil {
			return fmt.Sprintf(ErrorValidationLength, validationKey, length)
		}
	}

	return ""
}

func (v *Validation) validateTime(value string, validationKey string) string {
	if value != "" {
		_, err := time.Parse(time.RFC3339, value)

		if err != nil {
			return fmt.Sprintf(ErrorValidationTime, validationKey)
		}
	}

	return ""
}
