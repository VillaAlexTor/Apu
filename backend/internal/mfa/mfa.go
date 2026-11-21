package mfa

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1" // SHA-1 is required by RFC 6238's interoperable TOTP profile.
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	secretBytes       = 20
	recoveryCodeCount = 10
	recoveryCodeBytes = 10
	periodSeconds     = 30
)

var ErrInvalidKey = errors.New("MFA encryption key must contain 32 bytes")

type Encryptor struct {
	aead cipher.AEAD
}

func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create MFA cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create MFA GCM: %w", err)
	}
	return &Encryptor{aead: aead}, nil
}

func DecodeEncryptionKey(value string) ([]byte, error) {
	key, err := base64.RawStdEncoding.DecodeString(strings.TrimSpace(value))
	if err != nil || len(key) != 32 {
		return nil, ErrInvalidKey
	}
	return key, nil
}

func (encryptor *Encryptor) Encrypt(plaintext string) ([]byte, error) {
	nonce := make([]byte, encryptor.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("generate MFA nonce: %w", err)
	}
	return encryptor.aead.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

func (encryptor *Encryptor) Decrypt(ciphertext []byte) (string, error) {
	nonceSize := encryptor.aead.NonceSize()
	if len(ciphertext) <= nonceSize {
		return "", errors.New("invalid encrypted MFA secret")
	}
	plaintext, err := encryptor.aead.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
	if err != nil {
		return "", errors.New("decrypt MFA secret")
	}
	return string(plaintext), nil
}

func GenerateSecret() (string, error) {
	raw := make([]byte, secretBytes)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate MFA secret: %w", err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw), nil
}

func ProvisioningURI(issuer, account, secret string) string {
	label := url.PathEscape(issuer + ":" + account)
	values := url.Values{
		"algorithm": {"SHA1"},
		"digits":    {"6"},
		"issuer":    {issuer},
		"period":    {strconv.Itoa(periodSeconds)},
		"secret":    {secret},
	}
	return "otpauth://totp/" + label + "?" + values.Encode()
}

func Verify(secret, code string, now time.Time) bool {
	_, valid := VerifyStep(secret, code, now)
	return valid
}

func VerifyStep(secret, code string, now time.Time) (int64, bool) {
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return 0, false
	}
	for _, character := range code {
		if character < '0' || character > '9' {
			return 0, false
		}
	}
	counter := now.UTC().Unix() / periodSeconds
	for offset := int64(-1); offset <= 1; offset++ {
		step := counter + offset
		if step >= 0 && hmac.Equal([]byte(code), []byte(codeForCounter(secret, uint64(step)))) {
			return step, true
		}
	}
	return 0, false
}

func codeForCounter(secret string, counter uint64) string {
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return ""
	}
	message := make([]byte, 8)
	binary.BigEndian.PutUint64(message, counter)
	digest := hmac.New(sha1.New, decoded)
	_, _ = digest.Write(message)
	sum := digest.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	binaryCode := (uint32(sum[offset])&0x7f)<<24 |
		(uint32(sum[offset+1])&0xff)<<16 |
		(uint32(sum[offset+2])&0xff)<<8 |
		(uint32(sum[offset+3]) & 0xff)
	return fmt.Sprintf("%06d", binaryCode%1_000_000)
}

func GenerateRecoveryCodes() ([]string, [][]byte, error) {
	codes := make([]string, 0, recoveryCodeCount)
	hashes := make([][]byte, 0, recoveryCodeCount)
	for range recoveryCodeCount {
		raw := make([]byte, recoveryCodeBytes)
		if _, err := rand.Read(raw); err != nil {
			return nil, nil, fmt.Errorf("generate recovery code: %w", err)
		}
		encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw)
		code := encoded[:4] + "-" + encoded[4:8] + "-" + encoded[8:12] + "-" + encoded[12:16]
		codes = append(codes, code)
		hash := HashRecoveryCode(code)
		hashes = append(hashes, hash)
	}
	return codes, hashes, nil
}

func HashRecoveryCode(code string) []byte {
	normalized := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(code), "-", ""))
	digest := sha256.Sum256([]byte(normalized))
	return digest[:]
}
