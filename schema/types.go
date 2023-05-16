package schema

type Type struct {
	Name    string
	PkgPath string
	Kind    string

	// for Structs
	Fields []Field

	// for Funcs
	IsVariadic bool
	Ins        []Argument
	Outs       []Type
	Self       *Type // for Methods

	// for Maps
	Key *Type

	// for Arrays
	Len int

	// for Array,Chan,Map,Pointer,Slice
	Elem *Type
}

type Field struct {
	Name      string
	Type      Type
	Offset    uint
	Anonymous bool
}

type Argument struct {
	Name string
	Type Type
}

type Schema struct {
	All   []Type
	Types map[string]*Type
}
