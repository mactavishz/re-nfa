package main

import (
	"testing"

	"github.com/mactavishz/re-nfa/pkg/utils"
)

func TestNFAMatcher(t *testing.T) {
	tests := []struct {
		name     string
		regex    string
		input    string
		expected bool
	}{
		{"Simple match", "a", "a", true},
		{"Simple non-match", "a", "b", false},
		{"Star operator", "a*", "aaa", true},
		{"Star operator empty string", "a*", "", true},
		{"Plus operator", "a+", "aaa", true},
		{"Plus operator non-match", "a+", "", false},
		{"Optional operator", "a?b", "ab", true},
		{"Optional operator skipped", "a?b", "b", true},
		{"Alternation", "a|b", "a", true},
		{"Alternation other option", "a|b", "b", true},
		{"Alternation non-match", "a|b", "c", false},
		{"Grouping", "(ab)+", "abab", true},
		{"Grouping non-match", "(ab)+", "aba", false},
		{"Complex regex 1", "a(b|c)*d", "abcbcd", true},
		{"Complex regex 1 non-match", "a(b|c)*d", "abcbcde", false},
		{"Complex regex 2", "(a|b)c?d+", "acddd", true},
		{"Complex regex 2 non-match", "(a|b)c?d+", "ac", false},
		{"Complex regex 3", "(00)*|(000)*", "00", true},
		{"Complex regex 3", "(00)*|(000)*", "000", true},
		{"Complex regex 3 non-match", "(00)*|(000)*", "00000", false},
		{"Complex regex 3 non-match", "(00)*|(000)*", "0000000", false},
		{"Nested groups", "(a(b|c))+d", "abacd", true},
		{"Nested groups non-match", "(a(b|c))+d", "ababc", false},
		{"Nested alternation", "(a|b|c)(d|e|f)(g|h|i)", "bfh", true},
		{"Nested alternation non-match", "(a|b|c)(d|e|f)(g|h|i)", "bfj", false},
		{"Nested repetition", "(a+b*c?)+d", "aabcabd", true},
		{"Nested repetition non-match", "(a+b*c?)+d", "aabcab", false},
		{"Multiple group alternation", "(ab|cd)+(ef|gh)+", "abcdefghef", true},
		{"Multiple group alternation non-match", "(ab|cd)+(ef|gh)+", "abcdefg", false},
		{"Multiple alternations", "a|b|c|d", "c", true},
		{"Multiple alternations non-match", "a|b|c|d", "e", false},
		{"Combination of operators", "a+b*c?d", "aabd", true},
		{"Combination of operators non-match", "a+b*c?d", "bcd", false},
		{"Unicode support", "世界", "世界", true},
		{"Unicode support", "世界(和平)?", "世界和平", true},
		{"Unicode support non-match", "世界", "世界!", false},
		{"Unicode support non-match", "世界(和平)?", "世界不太和平", false},
		{"Unicode characters", "α+β+γ+", "ααβββγ", true},
		{"Unicode characters non-match", "α+β+γ+", "ααβββδ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := utils.NewParser(tt.regex)
			nfa, err := parser.Parse()
			if err != nil {
				t.Fatalf("Failed to parse regex '%s': %v", tt.regex, err)
			}

			result := nfa.Match(tt.input)
			if result != tt.expected {
				t.Errorf("Regex '%s' with input '%s': got %v, want %v", tt.regex, tt.input, result, tt.expected)
			}
		})
	}
}

func TestInvalidRegex(t *testing.T) {
	invalidRegexes := []string{
		"",
		"(",
		")",
		"*",
		"+",
		"?",
		"|",
		"(a|b",
		"a|b)",
		"a**",
		"a??",
		"a++",
		"()",
	}

	for _, regex := range invalidRegexes {
		t.Run(regex, func(t *testing.T) {
			parser := utils.NewParser(regex)
			_, err := parser.Parse()
			if err == nil {
				t.Errorf("Expected error for invalid regex '%s', but got nil", regex)
			}
		})
	}
}
