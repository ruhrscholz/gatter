package password

import "log"

type Password interface {
	DbEncode() string
	ValidateEncoded() bool
}

func GenerateFromPlaintext(alg string, password string) string {
	switch alg {
	case "argon2":
		return argon2GenerateFromPlaintext(password)
	default:
		log.Panicf("Unknown password algorithm %s", alg)
		return ""
	}
}
