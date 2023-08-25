package shortner

import (
	"testing"
	"time"
)

func TestShortner(t *testing.T) {
	urlschan := make(chan string)
	defer close(urlschan)

	shortner, err := New("testNode", urlschan)
	if err != nil {
		t.Fatalf("Failed to create Shortner: %v", err)
	}

	// Run the Short() method in a separate goroutine
	go shortner.Short()

	// Wait for some time to allow Shortner to generate URLs
	time.Sleep(1 * time.Second)

	select {
	case url := <-urlschan:
		if len(url) == 0 {
			t.Error("Generated URL should not be empty")
		}
	case <-time.After(20 * time.Microsecond):
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
