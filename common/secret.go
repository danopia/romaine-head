package common

import (
  "fmt"
  "encoding/base64"
  "crypto/rand"
)

func GenerateSecret() (string) {
  bytes := make([]byte, 32)

  _, err := rand.Read(bytes)
  if err != nil {
    fmt.Println("secret generation error:", err)
    return "error" // TODO
  }

  return base64.StdEncoding.EncodeToString(bytes)
}
