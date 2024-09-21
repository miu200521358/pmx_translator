package usecase

import (
	"testing"
)

func TestIsJapaneseString(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"こんにちは", true},      // Hiragana
		{"コンニチハ", true},      // Katakana
		{"ｺﾝﾆﾁﾊ", true},      // Half-width Katakana
		{"今日は", true},        // Kanji
		{"Hello", true},      // English
		{"こんにちはHello", true}, // Mixed Japanese and English
		{"こんにちは123", true},   // Mixed Japanese and numbers
		{"", true},           // Empty string
		{" ", true},          // Space
		{"\n", true},         // Newline
		{"こんにちは@", true},     // Prohibited symbol
		{"头饰4", false},       // Chinese
		{"こんにちは头饰4", false},  // Mixed Japanese and Chinese
		{"こんにちは🍣", false},    // Emoji
		{"こんにちは🍣🍣", false},   // Mixed Japanese and Emoji
	}

	ks, _ := loadKanji()

	for _, test := range tests {
		result := isJapaneseString(ks, test.input)
		if result != test.expected {
			t.Errorf("isJapaneseString(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
