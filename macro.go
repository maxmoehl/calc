package calc

import (
	"github.com/maxmoehl/calc/types"
)

// macroIndex maps the identifier of a macro to a function that can be used to create
// a macro for that identifier.
var macroIndex map[string]types.NewMacro

type macro struct {
	m types.Macro
}

func (m *macro) Locked() bool {
	return true
}

func (m *macro) Eval() (float64, error) {
	return m.m.Eval()
}