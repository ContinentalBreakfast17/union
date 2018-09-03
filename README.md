use 'godoc cmd/union' for documentation on the union command 

PACKAGE DOCUMENTATION

package union
    import "union"

    Package union provides an implementation for C-like unions

TYPES

type Union struct {
    // contains filtered or unexported fields
}
    Union is a C-like union containing internal data storage and a copy of
    the struct used to initialize it

func NewUnion(strct interface{}) *Union
    NewUnion creates a new union using the maximum sized field in the
    struct. If the given interface is not a struct, the function will return
    nil.

    The union will work on simple data types such as numbers. The union can
    handle 1 layer of nested structs also composed of simple data types.
    Other data types such as pointers, slices, channels, functions, arrays,
    maps, and so on are untested and probably will not work that well.
    Nesting additional structs will probably not work either.

    Nesting structs:

    The following should work: type nested struct {

	i int32
	f float64

    }

    type unionable struct {

	c complex 128
	s nested

    }

    However, adding additional structs to the definition of 'nested' will
    likely not work well.

func (u *Union) Get(f string) interface{}

	Get returns an interface containing the value associated with the
    field specified by f. This can be assumed to have the correct type-- no
    extra error checking should be necessary. f should be the name of the
    field in the struct definition, case sensitive. If f is not a valid
    field, Get returns nil.

func (u *Union) Set(f string, i interface{}) error
    Set sets the union data according the the field type and data specified.
    Like in Get, f should be the name of the field in the struct definition.
    The value of i must match the type of the field. If f is invalid or the
    type does not match the type of i, an error will be returned.
    Additionally, Set will panic if called on an unexported struct field.


