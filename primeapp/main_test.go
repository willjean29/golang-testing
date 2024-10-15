package main

import (
	"io"
	"os"
	"testing"
)

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"Zero", 0, false, "0 is not prime, by definition!"},
		{"One", 1, false, "1 is not prime, by definition!"},
		{"Negative", -1, false, "Negative numbers are not prime, by definition!"},
		{"Two", 2, true, "2 is a prime number!"},
		{"Three", 3, true, "3 is a prime number!"},
		{"Fifteen", 15, false, "15 is not prime, because it is a factor of 3"},
		{"SevenTeen", 17, true, "17 is a prime number!"},
	}

	for _, tt := range primeTests {
		t.Run(tt.name, func(t *testing.T) {
			result, msg := isPrime(tt.testNum)

			if result != tt.expected {
				t.Errorf("isPrime(%d) = %v; want %v", tt.testNum, result, tt.expected)
			}

			if msg != tt.msg {
				t.Errorf("isPrime(%d) = %s; want %s", tt.testNum, msg, tt.msg)
			}
		})
	}
}

func Test_prompt(t *testing.T) {
	// Redirect stdout to a pipe
	oldOut := os.Stdout
	// Create a pipe for stdout
	r, w, _ := os.Pipe()
	// set os.Stdout to the write end of the pipe
	os.Stdout = w

	prompt()

	// Close the write end of the pipe
	_ = w.Close()

	// Reset os.Stdout to its original value
	os.Stdout = oldOut

	// Read the output from the read end of the pipe
	out, _ := io.ReadAll(r)

	if string(out) != "->\n" {
		t.Errorf("prompt() = %s; want ->", out)
	}
}
