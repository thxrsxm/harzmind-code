package internal

import (
	"testing"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"Valid HTTP URL", "https://www.example.com", true},
		{"Valid HTTPS URL", "https://example.com/path/to/resource", true},
		{"Invalid URL (missing protocol)", "www.example.com", false},
		{"Invalid URL (invalid protocol)", "ftp://example.com", false},
		{"Localhost URL", "http://localhost:8080", true},
		{"Localhost URL with path", "http://localhost:8080/path", true},
		{"Empty URL", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidURL(tt.url); got != tt.want {
				t.Errorf("IsValidURL(%s) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}
