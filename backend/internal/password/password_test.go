package password

import "testing"

func TestHashAndVerify(t *testing.T) {
	hash, err := Hash("a-secure-local-password")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	valid, err := Verify("a-secure-local-password", hash)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !valid {
		t.Fatal("expected password to be valid")
	}

	valid, err = Verify("different-password", hash)
	if err != nil {
		t.Fatalf("Verify mismatch: %v", err)
	}
	if valid {
		t.Fatal("expected password mismatch")
	}
}

func TestHashRejectsShortPassword(t *testing.T) {
	if _, err := Hash("too-short"); err == nil {
		t.Fatal("expected short password to fail")
	}
}

func TestVerifyRejectsMalformedHash(t *testing.T) {
	if _, err := Verify("password", "not-a-hash"); err == nil {
		t.Fatal("expected malformed hash to fail")
	}
}
