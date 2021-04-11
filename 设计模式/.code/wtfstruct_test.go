package singleton

import "testing"

func TestStructA(t *testing.T) {
	var sa1, sa2 struct{}
	if sa1 == sa2 {
		t.Logf(`
Why GetInstance1() and GetInstance2() may return the variables with the same address.
See https://golang.org/ref/spec#Size_and_alignment_guarantees
"A struct or array type has size zero if it contains no fields (or elements, respectively) that have a size greater than zero. Two distinct zero-size variables may have the same address in memory."
More details: https://stackoverflow.com/questions/48052722/addresses-of-slices-of-empty-structs`)
	}
}
