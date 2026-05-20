package model

import "testing"

func TestValidateValidProduct(t *testing.T) {
	p := Product{
		Name:  "Test Product",
		Price: 10.99,
	}

	if !p.Validate() {
		t.Error("expected product to be valid")
	}
}

func TestValidateEmptyName(t *testing.T) {
	p := Product{
		Name:  "",
		Price: 10.99,
	}

	if p.Validate() {
		t.Error("expected product with empty name to be invalid")
	}
}

func TestValidateNegativePrice(t *testing.T) {
	p := Product{
		Name:  "Test Product",
		Price: -1,
	}

	if p.Validate() {
		t.Error("expected product with negative price to be invalid")
	}
}

func TestValidateZeroPrice(t *testing.T) {
	p := Product{
		Name:  "Free Product",
		Price: 0,
	}

	if !p.Validate() {
		t.Error("expected zero price to be valid")
	}
}

func TestValidateEmptyNameAndNegativePrice(t *testing.T) {
	p := Product{
		Name:  "",
		Price: -5,
	}

	if p.Validate() {
		t.Error("expected product with empty name and negative price to be invalid")
	}
}
