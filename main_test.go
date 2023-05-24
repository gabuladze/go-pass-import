package main

import (
	"testing"
)

func TestFormatPassword(t *testing.T) {
	testRow := []string{"https://test.test.com/test", "TestUsername", "testPaSsWord", "", "", "Test", "Test", "false"}
	expectedOutput := "testPaSsWord\nURL: https://*.test.com/*\nUsername: TestUsername\nExtra:\n\n"

	result := formatPassword(testRow)
	if result != expectedOutput {
		t.Errorf("formatPassword(testRow) got: %v, expected: %v", result, expectedOutput)
	}
}
