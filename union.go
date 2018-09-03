// Package union provides an implementation for C-like unions
package union

import (
	"errors"
	"reflect"
	"unsafe"
)

// Union is a C-like union containing internal data storage and a copy of
// the struct used to initialize it
type Union struct {
	data 	[]byte
	strct 	interface{}
}

// NewUnion creates a new union using the maximum sized field in the struct. 
// If the given interface is not a struct, the function will return nil.
//
// The union will work on simple data types such as numbers. The union
// can handle 1 layer of nested structs also composed of simple
// data types. Other data types such as pointers, slices, channels,
// functions, arrays, maps, and so on are untested and probably
// will not work that well. Nesting additional structs will probably
// not work either.
//
// Nesting structs:
//
// The following should work:
// type nested struct {
//     i int32
//     f float64
// }
//
// type unionable struct {
//     c complex 128
//     s nested
// }
//
// However, adding additional structs to the definition of 'nested'
// will likely not work well.
func NewUnion(strct interface{}) (*Union) {
	t := reflect.TypeOf(strct)
	if t.Kind() != reflect.Struct {
		return nil
	}

	maxSize := uintptr(0)
	for i := 0; i < t.NumField(); i++ {
		if size := t.Field(i).Type.Size(); size > maxSize {
			maxSize = size
		}
	}

	return &Union{make([]byte, maxSize, maxSize), strct}
}

//  Get returns an interface containing the value associated with the
// field specified by f. This can be assumed to have the correct
// type-- no extra error checking should be necessary. f should
// be the name of the field in the struct definition, case sensitive.
// If f is not a valid field, Get returns nil.
func (u *Union) Get(f string) interface{} {
	v := reflect.ValueOf(u.strct)
	field := v.FieldByName(f)
	if !field.IsValid() {
		return nil
	}
	
	ptr := reflect.NewAt(field.Type(), unsafe.Pointer(&u.data[0]))
	return reflect.Indirect(ptr)
}

// Set sets the union data according the the field type and data specified.
// Like in Get, f should be the name of the field in the struct definition.
// The value of i must match the type of the field. If f is invalid
// or the type does not match the type of i, an error will be returned.
// Additionally, Set will panic if called on an unexported struct field.
func (u *Union) Set(f string, i interface{}) error {
	tmp := reflect.ValueOf(u.strct)
	s := reflect.New(tmp.Type()).Elem()
	v := reflect.ValueOf(i)

	field := s.FieldByName(f)
	if !field.IsValid() {
		return errors.New("Invalid field")
	} else if v.Type().Kind() != field.Type().Kind() {
		return errors.New("Mismatched types")
	}

	field.Set(reflect.Indirect(v))
	addr := unsafe.Pointer(field.UnsafeAddr())

	for i := uintptr(0); i < uintptr(len(u.data)); i++ {
		if i < v.Type().Size() {
			b := unsafe.Pointer(uintptr(addr) + i*unsafe.Sizeof(u.data[0]))
			u.data[i] = *(*byte)(b)
		} else {
			u.data[i] = 0
		}
	}
	
	return nil
}