package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	slackApiVersion = "v0"
)

func BuildMessageForAuth(body, timestamp string, key []byte) []byte {
	return []byte(BuildMessageForAuthString(body, timestamp, key))
}

func BuildMessageForAuthString(body, timestamp string, key []byte) string {
	message := fmt.Sprintf("%s:%s:%s", slackApiVersion, timestamp, body)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return fmt.Sprintf("v0=%s", hex.EncodeToString(mac.Sum(nil)))
}
