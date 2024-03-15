package middlewares

import (
	"encoding/base64"
	"strings"
)

func SliceCommand(command string) []string {
	slicedCommand := strings.Split(strings.TrimSuffix(command, "\n"), " ")
	return slicedCommand
}

func B64Encode(resp string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(resp))
	return b64
}

