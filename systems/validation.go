package systems

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	validator "gopkg.in/go-playground/validator.v8"
)

type Validation struct {
	ErrorMessages map[string]string
	Validator     *validator.Validate
}

func (v *Validation) Validate(values map[string][]string, validationTags map[string]string) *ErrorData {
	v.initValidator()

	for validationKey, validationRules := range validationTags {

		validationRules := strings.Split(validationRules, ",")

		for _, validationRule := range validationRules {
			var message string

			validationRule := strings.Split(validationRule, "=")

			switch validationRule[0] {
			case "required":
				message = v.validateRequired(values[validationKey], validationKey)
			case "uuid5":
				message = v.validateUUIDV5(values[validationKey], validationKey)
			case "alpha":
				message = v.validateAlpha(values[validationKey], validationKey)
			case "alphanum":
				message = v.validateAlphaNumeric(values[validationKey], validationKey)
			case "numeric":
				message = v.validateNumeric(values[validationKey], validationKey)
			case "min":
				message = v.validateMin(values[validationKey], validationKey, validationRule[1])
			case "max":
				message = v.validateMax(values[validationKey], validationKey, validationRule[1])
			case "email":
				message = v.validateEmail(values[validationKey], validationKey)
			case "len":
				message = v.validateLength(values[validationKey], validationKey, validationRule[1])
			case "time":
				message = v.validateTime(values[validationKey], validationKey)
			case "latitude":
				message = v.validateLatitude(values[validationKey], validationKey)
			case "longitude":
				message = v.validateLongitude(values[validationKey], validationKey)
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

func (v *Validation) validateRequired(values []string, validationKey string) string {
	if len(values) == 0 {
		return fmt.Sprintf(ErrorValidationRequired, validationKey)
	}

	for _, value := range values {
		err := v.Validator.Field(value, "required")

		if err != nil {
			return fmt.Sprintf(ErrorValidationRequired, validationKey)
		}
	}

	return ""
}

func (v *Validation) validateUUIDV5(values []string, validationKey string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "uuid5")

			if err != nil {
				return fmt.Sprintf(ErrorValidationUUID5, validationKey)
			}
		}
	}

	return ""
}

func (v *Validation) validateAlpha(values []string, validationKey string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "alpha")

			if err != nil {
				return fmt.Sprintf(ErrorValidationAlpha, validationKey)
			}
		}
	}

	return ""
}

func (v *Validation) validateAlphaNumeric(values []string, validationKey string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "alphanum")

			if err != nil {
				return fmt.Sprintf(ErrorValidationAlphaNum, validationKey)
			}
		}
	}

	return ""
}

func (v *Validation) validateNumeric(values []string, validationKey string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "numeric")

			if err != nil {
				return fmt.Sprintf(ErrorValidationNumeric, validationKey)
			}
		}
	}

	return ""
}

func (v *Validation) validateMin(values []string, validationKey string, minLength string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "min"+"="+minLength)

			if err != nil {
				return fmt.Sprintf(ErrorValidationMin, validationKey, minLength)
			}
		}
	}

	return ""
}

func (v *Validation) validateMax(values []string, validationKey string, maxLength string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "max"+"="+maxLength)

			if err != nil {
				return fmt.Sprintf(ErrorValidationMax, validationKey, maxLength)
			}
		}
	}

	return ""
}

func (v *Validation) validateEmail(values []string, validationKey string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "email")

			if err != nil {
				return fmt.Sprintf(ErrorValidationEmail, validationKey)
			}
		}
	}

	return ""
}

func (v *Validation) validateLength(values []string, validationKey string, length string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "len"+"="+length)

			if err != nil {
				return fmt.Sprintf(ErrorValidationLength, validationKey, length)
			}
		}
	}

	return ""
}

func (v *Validation) validateTime(values []string, validationKey string) string {
	for _, value := range values {
		if value != "" {
			_, err := time.Parse(time.RFC3339, value)

			if err != nil {
				return fmt.Sprintf(ErrorValidationTime, validationKey)
			}
		}
	}

	return ""
}

func (v *Validation) validateLatitude(values []string, validationKey string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "latitude")

			if err != nil {
				return fmt.Sprintf(ErrorValidationLatitude, validationKey)
			}
		}
	}

	return ""
}

func (v *Validation) validateLongitude(values []string, validationKey string) string {
	for _, value := range values {
		if value != "" {
			err := v.Validator.Field(value, "longitude")

			if err != nil {
				return fmt.Sprintf(ErrorValidationLongitude, validationKey)
			}
		}
	}

	return ""
}
