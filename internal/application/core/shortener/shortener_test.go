package shortener

import (
	"errors"
	"testing"
	"time"
)

func TestShortener(t *testing.T) {
	urlschan := make(chan string)

	shortener, err := New("testNode", urlschan)
	if err != nil {
		t.Fatalf("Failed to create shortener: %v", err)
	}

	// Run the Short() method in a separate goroutine
	go shortener.Short()

	// Wait for some time to allow shortener to generate URLs
	time.Sleep(1 * time.Second)

	select {
	case url := <-urlschan:
		if len(url) == 0 {
			t.Error("Generated URL should not be empty")
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout: No URL generated within expected time")
	}
}

func TestDJB2(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"test1", 479},
		{"test2", 480},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := djb2(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected %d, but got %d", test.input, test.expected, result)
		}
	}
}

func TestIPLower(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		err      error
	}{
		{"192.168.1.123", 379, nil},
		{"10.0.0.255", 255, nil},
		{"invalidip", 0, errors.New("Invalid IP")},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", 0, errors.New("IPv6 address not supported")},
	}

	for _, test := range tests {
		result, err := iplower(test.input)

		if result != test.expected {
			t.Errorf("For input %s, expected %d, but got %d", test.input, test.expected, result)
		}

		if (err == nil && test.err != nil) || (err != nil && test.err == nil) || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
			t.Errorf("For input %s, expected error: %v, but got error: %v", test.input, test.err, err)
		}
	}
}
