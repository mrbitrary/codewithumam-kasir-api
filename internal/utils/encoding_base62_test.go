package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "encode UUID string",
			input: "550e8400-e29b-41d4-a716-446655440000",
			want:  "1g9Ej8dHSqGvKZMdHPjdgYCqLJ",
		},
		{
			name:  "encode simple string",
			input: "hello",
			want:  "5TP3P3v",
		},
		{
			name:  "encode empty string",
			input: "",
			want:  "0",
		},
		{
			name:  "encode single character",
			input: "a",
			want:  "2R",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeBase62(tt.input)
			assert.NotEmpty(t, got, "EncodeBase62 should not return empty string")
		})
	}
}

func TestDecodeBase62(t *testing.T) {
	// Test with actual encode/decode pairs
	original := "hello"
	encoded := EncodeBase62(original)
	decoded := DecodeBase62(encoded)
	
	assert.Equal(t, original, decoded, "Decode should reverse encode")
	
	// Test zero
	assert.Equal(t, "", DecodeBase62("0"))
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	tests := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"hello world",
		"test123",
		"a",
	}

	for _, original := range tests {
		t.Run(original, func(t *testing.T) {
			encoded := EncodeBase62(original)
			decoded := DecodeBase62(encoded)
			assert.Equal(t, original, decoded, "Round trip should preserve original value")
		})
	}
}
