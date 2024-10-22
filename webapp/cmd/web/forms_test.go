package main

import (
	"testing"
)

func TestForm_Has(t *testing.T) {
	test := []struct {
		name     string
		data     map[string][]string
		field    string
		expected bool
	}{
		{"field not found", map[string][]string{"name": {"John"}}, "email", false},
		{"field found", map[string][]string{"name": {"John"}}, "name", true},
	}

	for _, tt := range test {
		form := NewForm(tt.data)
		t.Run(tt.name, func(t *testing.T) {
			if form.Has(tt.field) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, !tt.expected)
			}
		})
	}
}

func TestForm_Required(t *testing.T) {
	test := []struct {
		name     string
		data     map[string][]string
		field    string
		expected string
	}{
		{"valid required values", map[string][]string{"name": {"John"}}, "name", ""},
		{"invalid required values", map[string][]string{"name": {"John"}}, "email", "This field cannot be blank"},
	}

	for _, tt := range test {
		form := NewForm(tt.data)
		t.Run(tt.name, func(t *testing.T) {
			form.Required(tt.field)
			if form.Errors.Get(tt.field) != tt.expected {
				t.Errorf("expected error message, got empty")
			}
		})
	}
}

func TestForm_check(t *testing.T) {
	form := NewForm(map[string][]string{"name": {"John"}})

	form.Check(false, "email", "This field cannot be blank")

	if len(form.Errors) == 0 {
		t.Errorf("expected error message, got empty")
	}
}

func TestForm_Valid(t *testing.T) {
	test := []struct {
		name     string
		data     map[string][]string
		field    string
		expected bool
	}{
		{"valid form", map[string][]string{"name": {"John"}}, "name", true},
		{"invalid form", map[string][]string{"name": {"John"}}, "email", false},
	}

	for _, tt := range test {
		form := NewForm(tt.data)
		t.Run(tt.name, func(t *testing.T) {
			form.Required(tt.field)
			if form.Valid() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, !tt.expected)
			}
		})
	}
}
