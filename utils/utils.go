package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"hash/crc32"

	"github.com/hashicorp/errwrap"
)

func GetRandomID() (string, error) {
	bytes := make([]byte, 10)
	n, err := rand.Reader.Read(bytes)
	if n != 10 {
		return "", errors.New("generated insufficient random bytes")
	}
	if err != nil {
		return "", errwrap.Wrapf("error generating random bytes: {{err}}", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func GetHash(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func GetSlice(slice []interface{}) (stringSLice []string) {
	for _, sl := range slice {
		stringSLice = append(stringSLice, sl.(string))
	}
	return
}
