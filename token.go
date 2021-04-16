package calc

// Token is an interface that gets generated by the lexer. A Token can either
// be a operator or a literal.
type Token interface {
	Type() string
	Value() interface{}
}

type token struct {
	tokenType string
	value     interface{}
}

func (t token) Type() string {
	return t.tokenType
}

func (t token) Value() interface{} {
	return t.value
}
