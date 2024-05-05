package utils

const (
	AttributeLocale       = "locale"
	AttributeName         = "name"
	AttributeReferralCode = "referral_code"
)

// init dictionary interface
func Dictionary() map[string]interface{} {
	res := make(map[string]interface{})
	return res
}
