package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"slices"

	terraformErrors "github.com/nikhilsbhat/terraform-provider-k3d/pkg/errors"
	"gopkg.in/yaml.v2"
)

// GetRandomID returns a random id when invoked.
func GetRandomID() (string, error) {
	randInt := 10
	bytes := make([]byte, randInt)

	n, err := rand.Reader.Read(bytes)
	if n != randInt {
		return "", terraformErrors.ErrInsufficientRandomBytes
	}

	if err != nil {
		return "", fmt.Errorf("%w: %w", terraformErrors.ErrGenerateRandomBytes, err)
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
func GetSlice(slice []any) []string {
	stringSLice := make([]string, 0, len(slice))
	for _, sl := range slice {
		stringSLice = append(stringSLice, sl.(string))
	}

	return stringSLice
}

// GetChecksum gets the checksum of passed string.
func GetChecksum(value string) (string, error) {
	cksm := sha256.New()

	if _, err := cksm.Write([]byte(value)); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(cksm.Sum(nil)), nil
}

// String returns string converted interface.
func String(value any) string {
	return value.(string)
}

// Bool returns bool converted interface.
func Bool(value any) bool {
	return value.(bool)
}

// Int returns bool converted interface.
func Int(value any) int {
	return value.(int)
}

// MapSlice returns array flattens the object passed to []map[string]any
// to simplify terraform attributes saving.
func MapSlice(value any) ([]map[string]any, error) {
	mp := make([]map[string]any, 0)

	j, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(j, &mp); err != nil {
		return nil, err
	}

	return mp, nil
}

// Map returns array flattens the object passed to []map[string]any
// to simplify terraform attributes saving.
func Map(value any) (map[string]string, error) {
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
func Yaml(data any) (string, error) {
	yml, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(yml), err
}

// JSON returns json encoded data structure passed to it.
func JSON(data any) (string, error) {
	jsn, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsn), err
}

// Contains returns true if given element is present the specified slice.
func Contains(s []string, searchTerm string) bool {
	return slices.Contains(s, searchTerm)
}
