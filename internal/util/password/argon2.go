package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func argon2GenerateFromPlaintext(password string) string {
	params := &argon2Params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	salt := make([]byte, params.saltLength)
	rand.Read(salt)

	hash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.memory,
		params.iterations,
		params.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawURLEncoding.EncodeToString(hash))
}

func argon2Validate(input string, correct string) bool {
	split := strings.Split(correct, "$")
	if len(split) != 6 {
		return false
	}

	params := argon2Params{}

	var version int
	_, err := fmt.Sscanf(split[2], "v=%d", version)
	if err != nil || version > argon2.Version {
		return false
	}

	_, err = fmt.Sscanf(split[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(split[4])
	if err != nil {
		return false
	}
	params.saltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(split[4])
	if err != nil {
		return false
	}
	params.keyLength = uint32(len(hash))

	inputHash := argon2.IDKey([]byte(input), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	if subtle.ConstantTimeCompare(inputHash, hash) == 1 {
		return true
	}
	return false
}
