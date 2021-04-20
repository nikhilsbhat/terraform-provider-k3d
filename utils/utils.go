package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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

func GetChecksum(value string) (string, error) {
	cksm := sha256.New()
	_, err := cksm.Write([]byte(value))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(cksm.Sum(nil)), nil
}

func String(value interface{}) string {
	return value.(string)
}

func Bool(value interface{}) bool {
	return value.(bool)
}

func Map(value interface{}) ([]map[string]interface{}, error) {
	mp := make([]map[string]interface{}, 0)
	j, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(j, &mp); err != nil {
		return nil, err
	}
	return mp, nil
}
