package keybackup

import "testing"

func TestValidatePayload(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		salt    string
		wantErr bool
	}{
		{"valid payload", "base64token", "salt", false},
		{"empty key", "", "salt", true},
		{"empty salt", "base64", "", true},
		{"whitespace only", "   ", "   ", true},
		{"too long", makeString(maxEncryptedPayload + 1), "salt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePayload(tt.key, tt.salt)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validatePayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func makeString(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = 'a'
	}
	return string(bytes)
}
