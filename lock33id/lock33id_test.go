package lock33id

import (
	"math"
	"testing"
)

func TestIdentityDerivations(t *testing.T) {
	x := math.Pi / 4
	witness, err := IdentityWitness(x)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(witness-math.Sqrt2/2) > 1e-12 {
		t.Fatalf("IdentityWitness(%v) = %v, want %v", x, witness, math.Sqrt2/2)
	}

	modulus, err := DeriveModulus(x)
	if err != nil {
		t.Fatal(err)
	}
	if modulus != 207 {
		t.Errorf("DeriveModulus(%v) = %d, want 207", x, modulus)
	}

	shift, err := DeriveShift(x)
	if err != nil {
		t.Fatal(err)
	}
	if shift != 25 {
		t.Errorf("DeriveShift(%v) = %d, want 25", x, shift)
	}

	seed, err := DeriveSeed(x)
	if err != nil {
		t.Fatal(err)
	}
	if seed != 7106 {
		t.Errorf("DeriveSeed(%v) = %d, want 7106", x, seed)
	}
}

func TestLock33PipeRoundTrip(t *testing.T) {
	const text = "Lock33 ☃"
	const wantPayload = "58|7b|6f|77|3f|3f|2c|260f"
	x := math.Pi / 4

	payload, err := EncodeLock33Pipe(text, x)
	if err != nil {
		t.Fatal(err)
	}
	if payload != wantPayload {
		t.Fatalf("EncodeLock33Pipe() = %q, want %q", payload, wantPayload)
	}

	decoded, err := DecodeLock33Pipe(payload, x)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != text {
		t.Errorf("DecodeLock33Pipe() = %q, want %q", decoded, text)
	}
}

func TestLock33PipeEmptyAndInvalidPayloads(t *testing.T) {
	encoded, err := EncodeLock33Pipe("", math.Pi/4)
	if err != nil {
		t.Fatal(err)
	}
	if encoded != "" {
		t.Errorf("EncodeLock33Pipe(empty) = %q, want empty", encoded)
	}
	decoded, err := DecodeLock33Pipe("", math.Pi/4)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != "" {
		t.Errorf("DecodeLock33Pipe(empty) = %q, want empty", decoded)
	}

	for _, payload := range []string{"1||2", "not-hex", "1"} {
		if _, err := DecodeLock33Pipe(payload, math.Pi/4); err == nil {
			t.Errorf("DecodeLock33Pipe(%q) succeeded, want an error", payload)
		}
	}
}

func TestIdentityRejectsUndefinedAngle(t *testing.T) {
	if _, err := IdentityWitness(math.Pi / 2); err == nil {
		t.Fatal("IdentityWitness(pi/2) succeeded, want an error")
	}
}
