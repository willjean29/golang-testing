package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

// nolint: funlen
func main() {
	runApp(os.Stdin, os.Stdout)
}

func runApp(input io.Reader, output io.Writer) {
	// print welcome message
	intro()

	// create a channel to communicate with the readUserInput goroutine
	doneChan := make(chan bool)

	// start the readUserInput goroutine
	go readUserInput(input, doneChan)

	// wait for the readUserInput goroutine to finish
	<-doneChan

	// close the doneChan channel
	close(doneChan)

	// print goodbye message
	fmt.Fprintln(output, "Goodbye!")
}

func intro() {
	fmt.Println("Welcome to the prime number app!")
	fmt.Println("This app will tell you if a number is prime or not.")
	fmt.Println("Let's get started!")
	prompt()
}

func prompt() {
	fmt.Println("->")
}

func readUserInput(io io.Reader, doneChan chan bool) {
	scanner := bufio.NewScanner(io)
	for {
		res, done := checkNumbers(scanner)

		if done {
			doneChan <- true
			return
		}
		fmt.Println(res)
		prompt()
	}
}

func checkNumbers(scanner *bufio.Scanner) (string, bool) {
	// read user input
	scanner.Scan()
	// check if user wants to quit
	if scanner.Text() == "q" {
		return "", true
	}

	//  try to convert user input to an integer
	numToCheck, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return "Please make a whole number!", false
	}

	// check if the number is prime
	_, msg := isPrime(numToCheck)

	return msg, false
}

func isPrime(n int) (bool, string) {
	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d is not prime, by definition!", n)
	}

	if n < 0 {
		return false, "Negative numbers are not prime, by definition!"
	}

	for i := 2; i < n/2; i++ {
		if n%i == 0 {
			return false, fmt.Sprintf("%d is not prime, because it is a factor of %d", n, i)
		}
	}

	return true, fmt.Sprintf("%d is a prime number!", n)
}
