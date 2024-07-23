package utils

import "testing"

func TestDecrypt(t *testing.T) {
	// Test cases
	tests := []struct {
		name         string
		encryptedHex string
		ivHex        string
		encodedKey   string
		expected     string
		expectError  bool
	}{
		{
			name:         "Invalid decryption",
			encryptedHex: "TBD",
			ivHex:        "TBD",
			encodedKey:   "TBD",
			expected:     "Hello, World!",
			expectError:  true,
		},
		{
			name:         "Invalid encoded key",
			encryptedHex: "f3ff7f2b5a69a3b1d454a12f38c40c29",
			ivHex:        "1234567890abcdef1234567890abcdef",
			encodedKey:   "invalid-base64",
			expected:     "",
			expectError:  true,
		},
		{
			name:         "Invalid encrypted hex",
			encryptedHex: "invalid-hex",
			ivHex:        "1234567890abcdef1234567890abcdef",
			encodedKey:   "MTIzNDU2Nzg5MDEyMzQ1Ng==",
			expected:     "",
			expectError:  true,
		},
		{
			name:         "Invalid IV hex",
			encryptedHex: "f3ff7f2b5a69a3b1d454a12f38c40c29",
			ivHex:        "invalid-hex",
			encodedKey:   "MTIzNDU2Nzg5MDEyMzQ1Ng==",
			expected:     "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decrypt(tt.encryptedHex, tt.ivHex, tt.encodedKey)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %q, but got %q", tt.expected, result)
				}
			}
		})
	}
}
