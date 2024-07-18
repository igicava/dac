package app

import (
	"testing"
)

func TestTokenizeAndInfinixToPostfix(t *testing.T) {
	// тестируем tokenize() 
	expression := "2.44 / 35 * 48 + 2 - 4.44"
	expectedTokens := []string{"2.44", "/", "35", "*", "48", "+", "2", "-", "4.44"}

	testTokenize := tokenize(expression)
	if len(expectedTokens) != len(testTokenize) {
		t.Fatalf("An unexpected number of items from tokenize. Want: %d, got: %d", len(expectedTokens), len(testTokenize))
	}
	for i := range expectedTokens {
		if expectedTokens[i] != testTokenize[i] {
			t.Fatalf("Invalid result: tokenize()")
		}
	}
	// тестируем infinixToPostfix
	expectedExpression := []string{"2.44", "35", "48", "*", "2", "+", "4.44", "/", "-"}
	if len(expectedExpression) != len(expectedTokens) {
		t.Fatalf("An unexpected number of items from infinixToPostfix. Want: %d, got: %d", len(expectedExpression), len(testTokenize))
	}
	for i := range expectedExpression {
		if expectedTokens[i] != testTokenize[i] {
			t.Fatalf("Invalid result: infinixToPostfix")
		}
	}
	
}