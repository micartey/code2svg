package code2svg

import (
	"encoding/base64"
	"strings"
)

func CalculateIndent(line string) (int, string) {
	trimmed := strings.TrimLeft(line, "\t ")
	prefix := line[:len(line)-len(trimmed)]
	tabs := strings.Count(prefix, "\t")
	spaces := strings.Count(prefix, " ")
	indent := tabs + (spaces / 4)
	return indent, trimmed
}

func DecodeBase64(input string) ([]byte, error) {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "+")

	// Try standard decoding
	if decoded, err := base64.StdEncoding.DecodeString(input); err == nil {
		return decoded, nil
	}
	// Try URL-safe decoding
	if decoded, err := base64.URLEncoding.DecodeString(input); err == nil {
		return decoded, nil
	}
	// Try raw (unpadded) decoding
	raw64 := strings.TrimRight(input, "=")
	if decoded, err := base64.RawStdEncoding.DecodeString(raw64); err == nil {
		return decoded, nil
	}
	// Final attempt: raw URL-safe
	return base64.RawURLEncoding.DecodeString(raw64)
}
