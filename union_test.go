package union

import (
	"fmt"
	"testing"
)

type sbuf [128]byte
type testUnion struct {
	u 	*Union
}

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
	A 	sbuf
	IA 	[8]int64
	s 	big
	P 	*int
}

func (t *testUnion) b() bool { return t.u.Get("b").(bool) }
func (t *testUnion) i() int64 { return t.u.Get("I").(int64) }
func (t *testUnion) f() float32 { return t.u.Get("f").(float32) }
func (t *testUnion) c() complex128 { return t.u.Get("c").(complex128) }
func (t *testUnion) a() sbuf { return t.u.Get("A").(sbuf) }
func (t *testUnion) ia() [8]int64 { return t.u.Get("IA").([8]int64) }
func (t *testUnion) s() big { return t.u.Get("s").(big) }
func (t *testUnion) p() *int { return t.u.Get("P").(*int) }

func (t testUnion) setI(v int64) { if err := t.u.Set("I", v); err != nil { panic(err) } }
func (t testUnion) setA(v sbuf) { if err := t.u.Set("A", v); err != nil { panic(err) } }
func (t testUnion) setIA(v [8]int64) { if err := t.u.Set("IA", v); err != nil { panic(err) } }

func TestUnion(t *testing.T) {
	union := &testUnion{NewUnion(test{})}
	if union.u == nil {
		panic("Failed to create union")
	}

	fmt.Println("Test 1:")
	union.setI(int64(10))
	fmt.Println(bufToS(union.a()))	
	fmt.Println(union.ia())

	fmt.Println("\nTest 2:")
	union.setA(toSBuf("\nhey"))
	fmt.Println(bufToS(union.a()))

	fmt.Println("\nTest 3:")



}

func toSBuf(s string) sbuf {
	var a sbuf
	copy(a[:], s)
	return a
}

func bufToS(b sbuf) string {
	return fmt.Sprintf("%s", b)
}