package utils

// TODO USE CONFIG TO STORE Array
var TestingDomains []string = []string{"localhost"}

func IsTestingDomain(domain string) bool {
	for _, testingDomain := range TestingDomains {
		if domain == testingDomain {
			return true
		}
	}
	return false
}
