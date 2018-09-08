// Package union provides an implementation for C-like unions
package union

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

/*
Union is a C-like union containing internal data storage and a copy of
the struct type used to initialize it
*/
type Union struct {
	data 	[]byte
	t 		reflect.Type
}

/*
NewUnion creates a new union using the maximum sized field in the struct. 
If the given interface is not a struct, the function will return an error.

The union has been confirmed to work on the following data types:
numbers, booleans, arrays, structs (can be nested), 
channels*, pointers*, functions*, maps*
* potential for runtime errors if you override these and then get them again

The following data types have finnicky behavior but will at least usually 
work on calls to set/get in succession:
slices, strings

Using interface{} in the union will not work well since the 
underlying type won't be interface{} when passing to Set
*/
func NewUnion(s interface{}) (*Union, error) {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("Value passed to NewUnion is not a struct")
	}

	maxSize := uintptr(0)
	for i := 0; i < t.NumField(); i++ {
		if size := t.Field(i).Type.Size(); size > maxSize {
			maxSize = size
		}
	}

	return &Union{make([]byte, maxSize, maxSize), t}, nil
}

/*
Get returns an interface containing the value associated with the
field specified by f. This can be assumed to have the correct
type-- no extra error checking should be necessary. f should
be the name of the field in the struct definition, case sensitive.
If f is not a valid field, Get returns nil.
*/
func (u *Union) Get(f string) interface{} {
	v := reflect.New(u.t).Elem()
	field := v.FieldByName(f)
	if !field.IsValid() {
		return nil
	}
	
	ptr := reflect.NewAt(field.Type(), unsafe.Pointer(&u.data[0]))
	return reflect.Indirect(ptr).Interface()
}

/*
Set sets the union data according the the field type and data specified.
Like in Get, f should be the name of the field in the struct definition.
The value of i must match the type of the field. If f is invalid
or the type does not match the type of i, Set will panic.
Additionally, Set will panic if called on an unexported struct field.
*/
func (u *Union) Set(f string, i interface{}) {
	s := reflect.New(u.t).Elem()
	v := reflect.ValueOf(i)

	field := s.FieldByName(f)
	if !field.IsValid() {
		panic(fmt.Sprintf("Union.Set called on non-existent field: %s", f))
	} else if v.Type().Kind() != field.Type().Kind() {
		panic(fmt.Sprintf("Union.Set cannot set field of type %s to value of type %s", field.Type(), v.Type().Kind()))
	}

	if !field.CanSet() {
		panic(fmt.Sprintf("Union.Set called on unexported field"))
	}
	field.Set(v)
	addr := unsafe.Pointer(field.UnsafeAddr())
	
	for i := uintptr(0); i < v.Type().Size(); i++ {
		b := unsafe.Pointer(uintptr(addr) + i*unsafe.Sizeof(u.data[0]))
		u.data[i] = *(*byte)(b)
	}	
}

/*
Wrap generates a wrapper for the union which can then be pasted into the code.
Using a wrapper will guarantee that union accesses do not have the wrong names.
*/
func (u *Union) Wrap() string {
	s := "type UnionWrapper struct {\n\tu *Union\n}\n"
	getters := ""
	setters := ""

	for i := 0; i < u.t.NumField(); i++ {
		field := u.t.Field(i)
		get := field.Name
		if field.PkgPath != "" {
			lower := strings.ToLower(field.Name)
			get = string(append([]byte{lower[0]}, []byte(get[1:])...))
		} else {
			setters += fmt.Sprintf("func (w *UnionWrapper) Set%s(v %s) { w.u.Set(\"%s\", v) }\n", field.Name, typeName(field), field.Name)
		}
		getters += fmt.Sprintf("func (w *UnionWrapper) %s() %s { return w.u.Get(\"%s\").(%s) }\n", get, typeName(field), field.Name, typeName(field))
	}

	return fmt.Sprintf("%s\n%s\n%s\n", s, getters, setters)
}

func typeName(field reflect.StructField) string {
	t := fmt.Sprintf("%s", field.Type)
	return strings.TrimPrefix(t, field.Type.PkgPath()+".")
}