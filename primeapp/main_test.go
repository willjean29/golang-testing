package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
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

func Test_intro(t *testing.T) {
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	intro()

	_ = w.Close()

	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)

	expected := "Welcome to the prime number app!\nThis app will tell you if a number is prime or not.\nLet's get started!\n->\n"

	if string(out) != expected {
		t.Errorf("intro() = %s; want %s", out, expected)
	}
}

func Test_checkNumbers(t *testing.T) {
	inputTests := []struct {
		name     string
		input    string
		expected string
		done     bool
	}{
		{"Quit", "q", "", true},
		{"Not Prime", "15", "15 is not prime, because it is a factor of 3", false},
		{"Prime", "7", "7 is a prime number!", false},
		{"Not a Number", "a", "Please make a whole number!", false},
	}
	for _, tt := range inputTests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			reader := bufio.NewScanner(input)
			res, _ := checkNumbers(reader)
			if res != tt.expected {
				t.Errorf("checkNumbers() = %s; want %s", res, tt.expected)
			}
		})
	}
}

func Test_readUserInput(t *testing.T) {
	doneChan := make(chan bool)
	input := "7\nq\n"
	var stdin bytes.Buffer
	stdin.Write([]byte(input))
	go readUserInput(&stdin, doneChan)
	<-doneChan
	close(doneChan)
}

func Test_runApp(t *testing.T) {
	input := "7\nq\n"
	var stdin bytes.Buffer
	stdin.Write([]byte(input))

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runApp(&stdin, w)

	os.Stdout = oldStdout
	_ = w.Close()

	out, _ := io.ReadAll(r)

	expected := "Welcome to the prime number app!\nThis app will tell you if a number is prime or not.\nLet's get started!\n->\n7 is a prime number!\n->\nGoodbye!\n"

	if string(out) != expected {
		t.Errorf("runApp() = %s; want %s", out, expected)
	}

}
