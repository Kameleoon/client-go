package utils

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"regexp"
	"strings"
)

const (
	HTTP         = "http://"
	HTTPS        = "https://"
	REGEX_DOMAIN = `^(\.?(([a-zA-Z\d][a-zA-Z\d-]*[a-zA-Z\d])|[a-zA-Z\d]))` +
		`(\.(([a-zA-Z\d][a-zA-Z\d-]*[a-zA-Z\d])|[a-zA-Z\d])){1,126}$`
	LOCALHOST = "localhost"
)

func ValidateTopLevelDomain(topLevelDomain string) string {
	if topLevelDomain == "" {
		return ""
	}

	topLevelDomain = strings.ToLower(topLevelDomain)

	protocols := []string{HTTP, HTTPS}
	for _, protocol := range protocols {
		if strings.HasPrefix(topLevelDomain, protocol) {
			topLevelDomain = strings.TrimPrefix(topLevelDomain, protocol)
			logging.Warning("The top-level domain contains %s. Domain after protocol trimming: %s", protocol, topLevelDomain)
			break
		}
	}

	matched, err := regexp.MatchString(REGEX_DOMAIN, topLevelDomain)
	if err != nil {
		logging.Error("Error compiling regex: %v", err)
		return topLevelDomain
	}

	if !matched && topLevelDomain != LOCALHOST {
		logging.Error(
			"The top-level domain %s is invalid. The value has been set as provided, but it does not meet"+
				" the required format for proper SDK functionality. Please check the domain for correctness.",
			topLevelDomain)
		return topLevelDomain
	}

	return topLevelDomain
}
