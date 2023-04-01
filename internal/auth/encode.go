package auth

import (
	b64 "encoding/base64"
	"idc-okta-api/internal/config"
)

func GetClientToken() string {
	data := config.GetOktaApiJksClientId() + ":" + config.GetOktaAPIClientSecret()
	retVal := b64.StdEncoding.EncodeToString([]byte(data))
	return retVal
}
