package common

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateSecret() string {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println("secret generation error:", err)
		return "error" // TODO
	}

	return base64.StdEncoding.EncodeToString(bytes)
}
