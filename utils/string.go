package utils

import (
	"strings"
	"encoding/hex"
	"crypto/rand"
)

func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	
	n, err := rand.Read(uuid)
	
	if n != len(uuid) || err != nil {
		return "", err
	}
	
	uuid[8] = 0x80
	uuid[4] = 0x40

	return hex.EncodeToString(uuid), nil
}

func EscapeForJSON(content string) (string) {
	myEscapedJSONString := strings.Replace(content, "\n", "\\n", -1)
	myEscapedJSONString = strings.Replace(myEscapedJSONString, "\"", "\\\"", -1)
	myEscapedJSONString = strings.Replace(myEscapedJSONString, "\r", "\\\r", -1)
	myEscapedJSONString = strings.Replace(myEscapedJSONString, "\t", "\\\t", -1)
	myEscapedJSONString = strings.Replace(myEscapedJSONString, "\b", "\\\b", -1)
	myEscapedJSONString = strings.Replace(myEscapedJSONString, "\f", "\\\f", -1)

	return myEscapedJSONString
}