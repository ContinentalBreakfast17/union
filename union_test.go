package union

import (
	"fmt"
	"testing"
)

type sbuf [128]byte

type big struct {
	i 	int32
	b 	bool
	c2 	complex128
}

type test struct {
	B 	bool
	I 	int64
	A 	sbuf
	IA 	[8]int64
	S 	big
	P 	*int
	Str string
	Ch 	chan bool
	Fnc func(i int) (string)
	M 	map[string]string
	f 	float32
}

type unionWrapper struct {
	u *Union
}

func (w *unionWrapper) B() bool { return w.u.Get("B").(bool) }
func (w *unionWrapper) I() int64 { return w.u.Get("I").(int64) }
func (w *unionWrapper) A() sbuf { return w.u.Get("A").(sbuf) }
func (w *unionWrapper) IA() [8]int64 { return w.u.Get("IA").([8]int64) }
func (w *unionWrapper) S() big { return w.u.Get("S").(big) }
func (w *unionWrapper) P() *int { return w.u.Get("P").(*int) }
func (w *unionWrapper) Str() string { return w.u.Get("Str").(string) }
func (w *unionWrapper) Ch() chan bool { return w.u.Get("Ch").(chan bool) }
func (w *unionWrapper) Fnc() func(int) string { return w.u.Get("Fnc").(func(int) string) }
func (w *unionWrapper) M() map[string]string { return w.u.Get("M").(map[string]string) }
func (w *unionWrapper) f() float32 { return w.u.Get("f").(float32) }

func (w *unionWrapper) SetB(v bool) { w.u.Set("B", v) }
func (w *unionWrapper) SetI(v int64) { w.u.Set("I", v) }
func (w *unionWrapper) SetA(v sbuf) { w.u.Set("A", v) }
func (w *unionWrapper) SetIA(v [8]int64) { w.u.Set("IA", v) }
func (w *unionWrapper) SetS(v big) { w.u.Set("S", v) }
func (w *unionWrapper) SetP(v *int) { w.u.Set("P", v) }
func (w *unionWrapper) SetStr(v string) { w.u.Set("Str", v) }
func (w *unionWrapper) SetCh(v chan bool) { w.u.Set("Ch", v) }
func (w *unionWrapper) SetFnc(v func(int) string) { w.u.Set("Fnc", v) }
func (w *unionWrapper) SetM(v map[string]string) { w.u.Set("M", v) }

func TestUnion(t *testing.T) {
	u, err := NewUnion(test{})
	if err != nil {
		panic(err)
	}
	fmt.Println(u.Wrap())
	union := &unionWrapper{u}

	fmt.Println("Test 1:")
	union.SetI(int64(21324569978))
	fmt.Println(bufToS(union.A()))	
	fmt.Println(union.IA())

	fmt.Println("\nTest 2:")
	union.SetA(toSBuf("hey"))
	union.SetB(false)
	fmt.Println(bufToS(union.A()))
	fmt.Println(union.I())

	fmt.Println("\nTest 3:")
	union.SetStr("reality")
	fmt.Println(union.Str())

	fmt.Println("\nTest 4:")
	i := 16
	union.SetP(&i)
	union.SetB(false)
	fmt.Println(*union.P()) // this isn't really safe...

	fmt.Println("\nTest 5:")
	ch := make(chan bool, 1)
	union.SetCh(ch)
	union.Ch() <- true
	fmt.Println(<-union.Ch())

	fmt.Println("\nTest 6:")
	union.SetFnc(testFunc)
	fmt.Println(union.Fnc()(17))

	fmt.Println("\nTest 7:")
	m := map[string]string{"poop": "descoop"}
	union.SetM(m)
	m["hey"] = "yo"
	union.M()["a"] = "b"
	fmt.Println(union.M())
	fmt.Println(union.f())

}

func toSBuf(s string) sbuf {
	var a sbuf
	copy(a[:], s)
	return a
}

func bufToS(b sbuf) string {
	return fmt.Sprintf("%s", b)
}

func testFunc(i int) string {
	return fmt.Sprintf("%d", i*2)
}