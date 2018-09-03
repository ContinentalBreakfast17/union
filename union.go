package main

import (
	"errors"
	"reflect"
	"unsafe"
)

type Union struct {
	data 	[]byte
	strct 	interface{}
}

func NewUnion(strct interface{}) (*Union) {
	t := reflect.TypeOf(strct)
	if t.Kind() != reflect.Struct {
		return nil
	}

	maxSize := uintptr(0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		kind := field.Type.Kind()
		size := field.Type.Size()

		switch kind {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.String, reflect.UnsafePointer:
			return nil
		}

		if size > maxSize {
			maxSize = size
		}
	}

	return &Union{make([]byte, maxSize, maxSize), strct}
}

func (u *Union) Get(typ string) interface{} {
	v := reflect.ValueOf(u.strct)
	field := v.FieldByName(typ)
	if !field.IsValid() {
		return nil
	}
	
	ptr := reflect.NewAt(field.Type(), unsafe.Pointer(&u.data[0]))
	return reflect.Indirect(ptr)
}

func (u *Union) Set(typ string, i interface{}) error {
	tmp := reflect.ValueOf(u.strct)
	s := reflect.New(tmp.Type()).Elem()
	v := reflect.ValueOf(i)

	field := s.FieldByName(typ)
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