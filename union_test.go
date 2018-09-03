package main

import (
	"fmt"
	"testing"
)

func TestUnion(t *testing.T) {
	type Big struct {
		c1 	complex128
		b 	bool
		c2 	complex128
	}

	type Test struct {
		b 	bool
		I 	int64
		f 	float64
		c	complex128
		s 	Big
	}

	
	union := NewUnion(Test{})
	if union == nil {
		return
	}

	err := union.Set("I", int64(456))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(union.Get("f"))
}