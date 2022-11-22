package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"

	"gopkg.in/yaml.v2"
)

// GetRandomID returns a random id when invoked.
func GetRandomID() (string, error) {
	randInt := 10
	bytes := make([]byte, randInt)
	n, err := rand.Reader.Read(bytes)
	if n != randInt {
		return "", errors.New("generated insufficient random bytes")
	}
	if err != nil {
		return "", fmt.Errorf("error generating random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// GetHash gets the hash of passed string.
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

// GetSlice returns StringSlice of passed interface array.
func GetSlice(slice []interface{}) []string {
	stringSLice := make([]string, 0)
	for _, sl := range slice {
		stringSLice = append(stringSLice, sl.(string))
	}

	return stringSLice
}

// GetChecksum gets the checksum of passed string.
func GetChecksum(value string) (string, error) {
	cksm := sha256.New()
	_, err := cksm.Write([]byte(value))
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(cksm.Sum(nil)), nil
}

// String returns string converted interface.
func String(value interface{}) string {
	return value.(string)
}

// Bool returns bool converted interface.
func Bool(value interface{}) bool {
	return value.(bool)
}

// Int returns bool converted interface.
func Int(value interface{}) int {
	return value.(int)
}

// MapSlice returns array flattens the object passed to []map[string]interface{}
// to simplify terraform attributes saving.
func MapSlice(value interface{}) ([]map[string]interface{}, error) {
	mp := make([]map[string]interface{}, 0)
	j, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(j, &mp); err != nil {
		return nil, err
	}

	return mp, nil
}

// Map returns array flattens the object passed to []map[string]interface{}
// to simplify terraform attributes saving.
func Map(value interface{}) (map[string]string, error) {
	var mp map[string]string
	j, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(j, &mp); err != nil {
		return nil, err
	}

	return mp, nil
}

// Encoder return base64 encoded string passed to it.
func Encoder(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// Yaml returns yaml encoded data structure passed to it.
func Yaml(data interface{}) (string, error) {
	yml, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(yml), err
}

// JSON returns json encoded data structure passed to it.
func JSON(data interface{}) (string, error) {
	jsn, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsn), err
}

// Contains returns true if given element is present the specified slice.
func Contains(s []string, searchTerm string) bool {
	for _, i := range s {
		if i == searchTerm {
			return true
		}
	}

	return false
}
