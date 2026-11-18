package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	defaultMemory  = 64 * 1024
	defaultTime    = 3
	defaultThreads = 4
	defaultKeyLen  = 32
	saltLength     = 16
)

var ErrInvalidHash = errors.New("invalid password hash")

func Hash(value string) (string, error) {
	if len(value) < 12 {
		return "", errors.New("password must contain at least 12 characters")
	}
	if len(value) > 1024 {
		return "", errors.New("password is too long")
	}

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	digest := argon2.IDKey([]byte(value), salt, defaultTime, defaultMemory, defaultThreads, defaultKeyLen)
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultMemory,
		defaultTime,
		defaultThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(digest),
	), nil
}

func Verify(value, encoded string) (bool, error) {
	parameters, salt, expected, err := parse(encoded)
	if err != nil {
		return false, err
	}

	actual := argon2.IDKey(
		[]byte(value),
		salt,
		parameters.time,
		parameters.memory,
		parameters.threads,
		uint32(len(expected)),
	)
	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}

type params struct {
	memory  uint32
	time    uint32
	threads uint8
}

func parse(encoded string) (params, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return params{}, nil, nil, ErrInvalidHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return params{}, nil, nil, ErrInvalidHash
	}

	parameterParts := strings.Split(parts[3], ",")
	if len(parameterParts) != 3 {
		return params{}, nil, nil, ErrInvalidHash
	}

	memory, err := parseUintParameter(parameterParts[0], "m=")
	if err != nil || memory == 0 || memory > 1024*1024 {
		return params{}, nil, nil, ErrInvalidHash
	}
	timeCost, err := parseUintParameter(parameterParts[1], "t=")
	if err != nil || timeCost == 0 || timeCost > 20 {
		return params{}, nil, nil, ErrInvalidHash
	}
	threadCount, err := parseUintParameter(parameterParts[2], "p=")
	if err != nil || threadCount == 0 || threadCount > 64 {
		return params{}, nil, nil, ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) < 16 {
		return params{}, nil, nil, ErrInvalidHash
	}
	digest, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(digest) < 16 || len(digest) > 64 {
		return params{}, nil, nil, ErrInvalidHash
	}

	return params{memory: uint32(memory), time: uint32(timeCost), threads: uint8(threadCount)}, salt, digest, nil
}

func parseUintParameter(value, prefix string) (uint64, error) {
	if !strings.HasPrefix(value, prefix) {
		return 0, ErrInvalidHash
	}
	return strconv.ParseUint(strings.TrimPrefix(value, prefix), 10, 32)
}
