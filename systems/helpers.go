package systems

import (
	"bytes"
	"fmt"
	"os"

	"github.com/elgs/gostrgen"
	uuid "github.com/satori/go.uuid"
	"github.com/ventu-io/go-shortid"
)

// Helpers Struct
type Helpers struct {
}

// StrConcat function used to concatenate multiple string
func (h *Helpers) StrConcat(args ...string) string {
	var buffer bytes.Buffer

	for index, element := range args {
		_ = index
		buffer.WriteString(element)
	}

	return buffer.String()
}

// GenerateUUID function use to generate Universal Unique Identifier using uuid v5
func (h *Helpers) GenerateUUID() string {
	namespaceDNS := uuid.NewV1()
	return uuid.NewV5(namespaceDNS, "api.shoppermate.com").String()
}

// GenerateUniqueShortID function used to generate unique short id.
func (h *Helpers) GenerateUniqueShortID() string {
	shortIDSettings, _ := shortid.New(1, shortid.DefaultABC, 123)

	shortID, _ := shortIDSettings.Generate()

	return shortID
}

// RandomString will generate random stringwith dynamic length and type of string
func (h *Helpers) RandomString(strType string, length int, include string, exclude string) string {
	// possible character sets are:
	// Lower, Upper, Digit, Punct, LowerUpper, LowerDigit, UpperDigit, LowerUpperDigit and All.
	// Any of the above can be combine by "|", e.g. LowerUpper is the same as Lower | Upper
	charSet := 0

	switch strType {
	case "Lower":
		charSet = gostrgen.Lower
	case "Upper":
		charSet = gostrgen.Upper
	case "Digit":
		charSet = gostrgen.Digit
	case "Punct":
		charSet = gostrgen.Punct
	case "LowerUpper":
		charSet = gostrgen.LowerUpper
	case "LowerDigit":
		charSet = gostrgen.LowerDigit
	case "UpperDigit":
		charSet = gostrgen.UpperDigit
	case "LowerUpperDigit":
		charSet = gostrgen.LowerUpperDigit
	default:
		charSet = gostrgen.All
	}

	charsToGenerate := length
	randomStr, err := gostrgen.RandGen(charsToGenerate, charSet, "", "")

	if err != nil {
		fmt.Println(err)
	}

	return randomStr

}

// StoragePath is a function to retrieve absolute storage path from config
// Return absolute storage path in string
func (h *Helpers) StoragePath() string {
	config := Configs{}
	return os.Getenv("GOPATH") + config.Get("app.yaml", "storage_path", "")
}

// StringInSlice function used to check if string contain in slice.
func (h *Helpers) StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func convertLogicalOperatorToSQLOperator(logicalOperator string) string {
	switch logicalOperator {
	case "eq":
		return "="
	case "ne":
		return "!="
	case "gt":
		return ">"
	case "ge":
		return ">="
	case "lt":
		return "<"
	case "le":
		return "<="
	case "and":
		return "AND"
	case "or":
		return "OR"
	}
	return "="
}
