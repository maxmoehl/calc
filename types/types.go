package types

// Node is the basic building block that the parser uses to build the abstract
// syntax tree.
type Node interface {
	// Locked indicates whether or not it is possible to modify this Node.
	// If false the parser will not try to modify this value. This currently
	// only returns false for the internal operation in certain cases.
	Locked() bool
	// Eval returns the value this Node resolves to, or an error if one occurs.
	Eval() (float64, error)
}

// Macro is the interface all macros have to implement.
type Macro interface {
	// Eval returns the value this macro resolves to, or an error if one occurs.
	Eval() (float64, error)
}

// NewMacro is a function a plugin needs to provide for every macro it contains.
// it is used to create a new in memory representation of the macro that can be
// used to evaluate it.
//
// The implementation should not evaluate any of its parameters but only check
// basic things like the number of parameters.
type NewMacro func(parameters []Node) (Macro, error)

// Index is a type that every plugin needs to have one instance of called `Index`.
// This map is used to get the identifiers of the macros and the names of their
// NewMacro functions.
type Index map[string]string
