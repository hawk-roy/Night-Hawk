package auth

import "testing"

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken(123, "jwt_user")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if claims.UserID != 123 {
		t.Fatalf("unexpected user_id: got %d want %d", claims.UserID, 123)
	}
	if claims.Username != "jwt_user" {
		t.Fatalf("unexpected username: got %q want %q", claims.Username, "jwt_user")
	}
}

func TestParseInvalidToken(t *testing.T) {
	_, err := ParseToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}
