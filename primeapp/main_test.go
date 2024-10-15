package main

import "testing"

func Test_isPrime(t *testing.T) {
	result, msg := isPrime(0)

	if result != false {
		t.Errorf("isPrime(0) = %v; want false", result)
	}

	if msg != "0 is not prime, by definition!" {
		t.Errorf("isPrime(0) = %v; want 0 is not prime, by definition!", msg)
	}
}
