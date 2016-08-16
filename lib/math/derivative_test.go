package invt_math

import (
	"testing"
)

func TestFirst(t *testing.T) {
	f := NewFirstDerivative()
	if f.Step(2) != 2 {
		t.Fatal("first step returned incorrect value")
	}
	if f.Step(3) != 1 {
		t.Fatal("second step returned incorrect value")
	}
	if f.Step(1) != -2 {
		t.Fatal("third step returned incorrect value")
	}
	if f.Step(4) != 3 {
		t.Fatal("fourth step returned incorrect value")
	}
}

func TestSecond(t *testing.T) {
	f := NewSecondDerivative()
	if f.Step(2) != 0 {
		t.Fatal("first step returned incorrect value")
	}
	x := f.Step(3)
	if x != -1 {
		t.Fatal(x)
		t.Fatal("second step returned incorrect value")
	}
	if f.Step(1) != -3 {
		t.Fatal("third step returned incorrect value")
	}
	if f.Step(4) != 5 {
		t.Fatal("fourth step returned incorrect value")
	}
}
