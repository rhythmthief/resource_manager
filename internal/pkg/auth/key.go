package auth

import (
	"crypto/rand"
	"encoding/base64"
)

/*GenerateKey creates a key of a given length and returns it to the caller
Args:	key length
Rets:	new key, error*/
func GenerateKey(length int) (string, error) {
	var err error
	var key string
	raw := make([]byte, length)

	//Populating raw with random bytes
	if _, err = rand.Read(raw); err == nil {
		//Encoding a URL-friendly string to be used as a key
		key = base64.URLEncoding.EncodeToString(raw)
	}

	return key, err
}
