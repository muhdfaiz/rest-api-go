package systems

import (
	"fmt"
	"net/http"
	"strconv"

	validator "gopkg.in/go-playground/validator.v8"
)

const (
	ValidationFailed             = "1001"
	FacebookIDNotValid           = "1002"
	DatabaseError                = "1003"
	ValueAlreadyExist            = "1004"
	BadRequest                   = "1005"
	InternalServerError          = "1006"
	BindingError                 = "1007"
	InvalidFileType              = "1008"
	FileSizeExceededLimit        = "1009"
	CannotReadFile               = "1010"
	CannotDeleteFile             = "1011"
	CannotCopyFile               = "1012"
	ErrorAmazonService           = "1013"
	ErrorConvertStringToInt      = "1014"
	ErrorConvertIntToString      = "1015"
	CannotDetectFileType         = "1016"
	FailedToSendSMS              = "1017"
	FailedToGenerateReferralCode = "1018"
	ReferralCodeNotExist         = "1019"
	ReferralCodeExceedLimit      = "1020"
	ResourceNotFound             = "1021"
	VerificationCodeInvalid      = "1022"
	CannotCreateResource         = "1023"
	FailedToGenerateToken        = "1024"
	TokenNotValid                = "1025"
	TokenIdentityNotMatch        = "1026"

	TitleValidationError         = "Validation failed."
	TitleInternalServerError     = "Internal server error."
	TitleFileSizeExceededLimit   = "File size exceeded the limit."
	TitleFileTypeError           = "Invalid file type."
	TitleDatabaseError           = "Database Error."
	TitleBindingError            = "Binding Error."
	TitleDuplicateValueError     = "%s already exists."
	TitleFacebookIDNotValidError = "Facebook ID not valid!"
	TitleFileEmptyError          = "File is empty."
	TitleSentSmsError            = "Sending Sms Failed."
	TitleReferralCodeNotExist    = "Referral code not exist."
	TitleReferralCodeExceedLimit = "Referral code Exceeded Limit."
	TitleResourceNotFoundError   = "%s not exists."
	TitleVerificationCodeInvalid = "Invalid verification code"
	TitleCannotCreateResouce     = "Failed to create new %s"
	TItleCannotUpdateResource    = "Failed to update %s with %s %s"
	TitleFailedToGenerateToken   = "Failed to generate Token"
	TitleErrorTokenNotValid      = "Access token error"
	TitleTokenIdentityNotMatch   = "Your access token belong to other user"

	ErrorValidationRequired = "The %s field is required."
	ErrorValidationUUID     = "The %s field is not valid uuid."
	ErrorValidationAlpha    = "The %s field may only contain letters."
	ErrorValidationAlphaNum = "The %s field may only contain letters and numbers."
	ErrorValidationNumeric  = "The %s field must be a number."
	ErrorValidationMin      = "The %s field length must be at least %s."
	ErrorValidationMax      = "The %s field length may not be greater than %s."
	ErrorValidationEmail    = "The %s field must be a valid email address."
	ErrorInternalServer     = "API cannot return results because an internal server error has occurred."

	ErrorFileSizeExceededLimit   = "The %s attributes may not be greater than %s kilobytes."
	ErrorFileType                = "You have uploaded an invalid file type. Only file type %s are allowed."
	ErrorDatabase                = "API cannot return results because an internal server error has occurred."
	ErrorBinding                 = "The request could not be understood by the server due to malformed syntax. Please check the request."
	ErrorDuplicateValue          = "Duplicate entry '%s' for key '%s'."
	ErrorFacebookIDNotValid      = "The facebook_id %s value is not valid."
	ErrorFileEmpty               = "File is empty."
	ErrorSentSms                 = "Please wait %s seconds before sending sms again."
	ErrorReferralCodeNotExist    = "Please enter correct referral code."
	ErrorReferralCodeExceedLimit = "The referral code you entered has exceeded the limit.Please use another referral code."
	ErrorResourceNotFound        = "%s with %s %s not exists in system."
	ErrorVerificationCodeInvalid = "The verification code you entered %s is invalid."
	ErrorTokenNotValid           = "The access token you sent could not be found or is invalid."
	ErrorTokenIdentityNotMatch   = "Cannot %s because your access token belong to other user. Please use your own access token."
)

type ErrorMessage struct{}

type Error struct{}

type ErrorData struct {
	Error *ErrorFormat `json:"errors"`
}

// ErrorFormat used to define structure for error message
type ErrorFormat struct {
	Status string      `json:"status"`
	Code   string      `json:"code"`
	Title  string      `json:"title"`
	Detail interface{} `json:"detail"`
}

// GenericError function used to create standardize custom error message
func (e Error) GenericError(status string, code string, title string, key string, value string) *ErrorData {
	if key == "" {
		key = "message"
	}
	return &ErrorData{
		Error: &ErrorFormat{
			Status: status,
			Code:   code,
			Title:  title,
			Detail: map[string]interface{}{key: value},
		},
	}
}

// InternalServerError function used to format error message when internal server error happen(500) happen
func (e Error) InternalServerError(errors interface{}, code string) *ErrorData {
	config := Configs{}
	errorFormat := &ErrorFormat{}
	errorFormat.Status = strconv.Itoa(http.StatusInternalServerError)
	errorFormat.Code = code
	errorFormat.Title = TitleInternalServerError
	errorFormat.Detail = errors

	if config.Get("app.yaml", "debug", "") == "false" || errors == nil {
		errorFormat.Detail = map[string]string{"message": ErrorInternalServer}
	}

	return &ErrorData{
		Error: errorFormat,
	}
}

// FileSizeExceededLimit used to format error message when file want to upload is bigger than size allowed by system
func (e Error) FileSizeExceededLimit(field string, maxFileSize string) *ErrorData {
	return &ErrorData{
		Error: &ErrorFormat{
			Status: strconv.Itoa(http.StatusRequestEntityTooLarge),
			Code:   FileSizeExceededLimit,
			Title:  TitleFileSizeExceededLimit,
			Detail: map[string]string{field: fmt.Sprintf(ErrorFileSizeExceededLimit, field, maxFileSize)},
		},
	}
}

// InvalidFileTypeError function used to format error message when file want to upload is not allowed
func (e Error) InvalidFileTypeError(allowFileType string) *ErrorData {
	return &ErrorData{
		Error: &ErrorFormat{
			Status: strconv.Itoa(http.StatusBadRequest),
			Code:   InvalidFileType,
			Title:  TitleFileTypeError,
			Detail: map[string]string{"message": fmt.Sprintf(ErrorFileType, allowFileType)},
		},
	}
}

func (e Error) TokenIdentityNotMatchError(text string) *ErrorData {
	return &ErrorData{
		Error: &ErrorFormat{
			Status: strconv.Itoa(http.StatusBadRequest),
			Code:   TokenIdentityNotMatch,
			Title:  TitleTokenIdentityNotMatch,
			Detail: map[string]interface{}{"message": fmt.Sprintf(ErrorTokenIdentityNotMatch, text)},
		},
	}
}

func (e Error) ResourceNotFoundError(resource string, attribute string, value string) *ErrorData {
	return &ErrorData{
		Error: &ErrorFormat{
			Status: strconv.Itoa(http.StatusBadRequest),
			Code:   ResourceNotFound,
			Title:  fmt.Sprintf(TitleResourceNotFoundError, resource),
			Detail: map[string]interface{}{attribute: fmt.Sprintf(ErrorResourceNotFound, resource, attribute, value)},
		},
	}
}

// DBError will return 500 sInternal Server Error
func (e Error) DBError(errors interface{}) *ErrorData {
	config := Configs{}
	errorFormat := &ErrorFormat{}
	errorFormat.Status = strconv.Itoa(http.StatusInternalServerError)
	errorFormat.Code = DatabaseError
	errorFormat.Title = TitleInternalServerError
	errorFormat.Detail = errors

	if config.Get("app.yaml", "debug", "") == "false" || errors == nil {
		errorFormat.Detail = map[string]string{"message": ErrorInternalServer}
	}

	return &ErrorData{
		Error: errorFormat,
	}
}

// BindingError will return 400 Bad Request Error
func (e Error) BindingError(errors interface{}) *ErrorData {
	config := Configs{}
	errorFormat := &ErrorFormat{}
	errorFormat.Status = strconv.Itoa(http.StatusBadRequest)
	errorFormat.Code = BadRequest
	errorFormat.Title = TitleBindingError
	errorFormat.Detail = errors

	if config.Get("app.yaml", "debug", "") == "false" || errors == nil {
		errorFormat.Detail = map[string]string{"message": ErrorBinding}
	}

	return &ErrorData{
		Error: errorFormat,
	}
}

// DuplicateValueErrors will handle 309 conflict error
func (e Error) DuplicateValueErrors(resourceType string, field string, value string) *ErrorData {
	return &ErrorData{
		Error: &ErrorFormat{
			Status: strconv.Itoa(http.StatusConflict),
			Code:   ValueAlreadyExist,
			Title:  fmt.Sprintf(TitleDuplicateValueError, resourceType),
			Detail: map[string]string{"message": fmt.Sprintf(ErrorDuplicateValue, field, value)},
		},
	}
}

// ValidationErrors will handle validation error messages
func (e Error) ValidationErrors(errors map[string]*validator.FieldError) *ErrorData {
	errorMessages := make(map[string]string)
	for _, errMsg := range errors {
		var message string

		// Set error message based on validation rule
		switch errMsg.ActualTag {
		case "required":
			message = fmt.Sprintf(ErrorValidationRequired, errMsg.Name)
		case "uuid4":
			message = fmt.Sprintf(ErrorValidationUUID, errMsg.Name)
		case "alpha":
			message = fmt.Sprintf(ErrorValidationAlpha, errMsg.Name)
		case "alphanum":
			message = fmt.Sprintf(ErrorValidationAlphaNum, errMsg.Name)
		case "numeric":
			message = fmt.Sprintf(ErrorValidationNumeric, errMsg.Name)
		case "min":
			message = fmt.Sprintf(ErrorValidationMin, errMsg.Name, errMsg.Param)
		case "max":
			message = fmt.Sprintf(ErrorValidationMax, errMsg.Name, errMsg.Param)
		case "email":
			message = fmt.Sprintf(ErrorValidationEmail, errMsg.Name)
		}
		errorMessages[errMsg.Name] = message
	}

	return &ErrorData{
		Error: &ErrorFormat{
			Status: strconv.Itoa(http.StatusBadRequest),
			Code:   ValidationFailed,
			Title:  TitleValidationError,
			Detail: errorMessages,
		},
	}
}
