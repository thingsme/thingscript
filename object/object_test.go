package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	origin1 := &Integer{Value: 1234567}
	origin2 := &Integer{Value: 1234567}
	diff1 := &Integer{Value: 3141592}
	diff2 := &Integer{Value: 3141592}

	if origin1.HashKey() != origin2.HashKey() {
		t.Errorf("integers with same value have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("integers with same value have different hash keys")
	}

	if origin1.HashKey() == diff1.HashKey() {
		t.Errorf("integers with different value have same hash keys")
	}
}

func TestFloatHashKey(t *testing.T) {
	origin1 := &Float{Value: 1.234567}
	origin2 := &Float{Value: 1.234567}
	diff1 := &Float{Value: 3.141592}
	diff2 := &Float{Value: 3.141592}

	if origin1.HashKey() != origin2.HashKey() {
		t.Errorf("floats with same value have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("floats with same value have different hash keys")
	}

	if origin1.HashKey() == diff1.HashKey() {
		t.Errorf("floats with different value have same hash keys")
	}
}
