package util

const (
	USD = "USD"
	EU  = "EU"
	CAD = "CAD"
	VND = "VND"
)

func IssupportedCurrency(currency string) bool {
	switch currency {
	case USD, EU, CAD, VND:
		return true
	}
	return false
}
