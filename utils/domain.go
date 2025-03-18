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

	topLevelDomain = checkAndTrimProtocol(strings.ToLower(topLevelDomain))

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

func ValidateNetworkDomain(networkDomain string) string {
	if networkDomain == "" {
		return ""
	}

	networkDomain = checkAndTrimProtocol(strings.ToLower(networkDomain))

	// replace first and last dot
	networkDomain = strings.Trim(networkDomain, ".")

	matched, err := regexp.MatchString(REGEX_DOMAIN, networkDomain)
	if err != nil {
		logging.Error("Error compiling regex: %v", err)
		return ""
	}

	if !matched && networkDomain != LOCALHOST {
		logging.Error("The network domain %s is invalid.", networkDomain)
		return ""
	}

	return networkDomain
}

func checkAndTrimProtocol(domain string) string {
	protocols := []string{HTTP, HTTPS}
	for _, protocol := range protocols {
		if strings.HasPrefix(domain, protocol) {
			domain = strings.TrimPrefix(domain, protocol)
			logging.Warning("The domain contains %s. Domain after protocol trimming: %s", protocol, domain)
			return domain
		}
	}
	return domain
}
