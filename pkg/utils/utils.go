package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// GetRandomID returns a random id when invoked.
func GetRandomID() (string, error) {
	randInt := 10
	bytes := make([]byte, randInt)
	n, err := rand.Reader.Read(bytes)
	if n != randInt {
		return "", fmt.Errorf("generated insufficient random bytes")
	}
	if err != nil {
		return "", fmt.Errorf("error generating random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
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

// String returns string converted interface.
func String(value interface{}) string {
	return value.(string)
}

// Bool returns bool converted interface.
func Bool(value interface{}) bool {
	return value.(bool)
}
