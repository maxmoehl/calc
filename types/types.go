package types

type Node interface {
	// Locked indicates whether or not it is possible to go down into the Operation.
	// If false, Operator, Left and Right might return unexpected values.
	Locked() bool
	// Eval returns the value this operation resolves to, or an error if one occurs.
	Eval() (float64, error)
}

type Operation interface {
	Operator() rune
	Left() Node
	Right() Node
	// Locked indicates whether or not it is possible to go down into the Operation.
	// If false, Operator, Left and Right might return unexpected values.
	Locked() bool
	// Eval returns the value this operation resolves to, or an error if one occurs.
	Eval() (float64, error)
}

type Macro interface {
	Eval() (float64, error)
}

type NewMacro func(parameter []Node) (Macro, error)

type Index map[string]string
