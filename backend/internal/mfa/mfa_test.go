package mfa

import (
	"encoding/base32"
	"testing"
	"time"
)

func TestTOTPMatchesRFCVector(t *testing.T) {
	secret := base64ToBase32(t, "12345678901234567890")
	if actual := codeForCounter(secret, 1); actual != "287082" {
		t.Fatalf("expected RFC code 287082, got %s", actual)
	}
	if !Verify(secret, "287082", time.Unix(59, 0)) {
		t.Fatal("expected RFC TOTP code to verify")
	}
}

func TestEncryptorRoundTrip(t *testing.T) {
	encryptor, err := NewEncryptor(make([]byte, 32))
	if err != nil {
		t.Fatalf("NewEncryptor: %v", err)
	}
	ciphertext, err := encryptor.Encrypt("JBSWY3DPEHPK3PXP")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	plaintext, err := encryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if plaintext != "JBSWY3DPEHPK3PXP" {
		t.Fatalf("unexpected plaintext %q", plaintext)
	}
}

func TestRecoveryCodesAreUnique(t *testing.T) {
	codes, hashes, err := GenerateRecoveryCodes()
	if err != nil {
		t.Fatalf("GenerateRecoveryCodes: %v", err)
	}
	if len(codes) != recoveryCodeCount || len(hashes) != recoveryCodeCount {
		t.Fatalf("expected %d recovery codes", recoveryCodeCount)
	}
	seen := map[string]bool{}
	for _, code := range codes {
		if seen[code] {
			t.Fatalf("duplicate recovery code %q", code)
		}
		seen[code] = true
	}
}

func base64ToBase32(t *testing.T, value string) string {
	t.Helper()
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(value))
}
