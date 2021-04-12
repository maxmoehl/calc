package main

type token interface {
	Type() string
	Value() interface{}
}

type operator struct {
	value rune
}

func (o operator) Type() string {
	return typeOperator
}

func (o operator) Value() interface{} {
	return o.value
}

type literal struct {
	value float64
}

func (l literal) Type() string {
	return typeLiteral
}

func (l literal) Value() interface{} {
	return l.value
}
