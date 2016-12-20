package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type ReferralCashbackRepositoryInterface interface {
	Create(referrerGUID string, referentGUID string) (interface{}, *systems.ErrorData)
	Count(conditionAttribute string, conditionValue string) int64
}
