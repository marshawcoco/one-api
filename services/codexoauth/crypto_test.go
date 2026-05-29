package codexoauth

import "testing"

func TestProtectRefreshTokenRoundTrip(t *testing.T) {
	protected := ProtectRefreshToken("refresh-token")
	if protected == "refresh-token" {
		t.Fatal("token should be protected")
	}
	unprotected, err := UnprotectRefreshToken(protected)
	if err != nil {
		t.Fatal(err)
	}
	if unprotected != "refresh-token" {
		t.Fatalf("token mismatch: %s", unprotected)
	}
}

func TestProtectRefreshTokenIsIdempotent(t *testing.T) {
	protected := ProtectRefreshToken("refresh-token")
	if ProtectRefreshToken(protected) != protected {
		t.Fatal("protect should be idempotent")
	}
}

func TestUnprotectRefreshTokenAllowsPlaintext(t *testing.T) {
	plain, err := UnprotectRefreshToken("refresh-token")
	if err != nil {
		t.Fatal(err)
	}
	if plain != "refresh-token" {
		t.Fatalf("token mismatch: %s", plain)
	}
}
