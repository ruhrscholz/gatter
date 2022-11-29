package password

import (
	"errors"
	"log"
	"strings"
)

func GenerateFromPlaintext(alg string, password string) (string, error) {
	switch alg {
	case "argon2":
		return argon2GenerateFromPlaintext(password), nil
	default:
		log.Panicf("Unknown password algorithm %s", alg)
		return "", errors.New("Unknown password algorithm")
	}
}

func Validate(input string, correct string) bool {
	switch {
	case strings.HasPrefix(correct, "$argon2"):
		return argon2Validate(input, correct)
	default:
		return false
	}
}
