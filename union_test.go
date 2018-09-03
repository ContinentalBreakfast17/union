package union

import (
	"fmt"
	"testing"
)

func TestUnion(t *testing.T) {
	type big struct {
		i 	int32
		b 	bool
		c2 	complex128
	}

	type test struct {
		b 	bool
		I 	int64
		f 	float32
		c	complex128
		s 	big
	}


	union := NewUnion(test{})
	if union == nil {
		return
	}

	err := union.Set("I", int64(456))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(union.Get("s"))
}
