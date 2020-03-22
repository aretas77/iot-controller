package utils

import "strings"

func StripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

func SplitTopic4(topic string) (string, string, string, string) {
	tok := strings.Split(topic, "/")
	return tok[0], tok[1], tok[2], tok[3]
}
