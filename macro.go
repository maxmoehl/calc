package calc

import (
	"fmt"
)

// macros maps the identifier of a macro to a function that can be used to create
// a macro for that identifier.
var macros map[string]func(nodes ...Operation) (Macro, error)

func RegisterMacro(identifier string, newMacroFunc func(nodes ...Operation) (Macro, error)) error {
	if macros == nil {
		macros = make(map[string]func(nodes ...Operation) (Macro, error))
	}
	if macros[identifier] != nil {
		return fmt.Errorf("identifier '%s' already in use", identifier)
	}
	macros[identifier] = newMacroFunc
	return nil
}

type Macro interface {
	Eval() (float64, error)
}

type macroOperation struct {
	m Macro
}

func (m *macroOperation) Operator() rune {
	return 'm'
}

func (m *macroOperation) Left() Operation {
	return nil
}

func (m *macroOperation) Right() Operation {
	return nil
}

func (m *macroOperation) Locked() bool {
	return true
}

func (m *macroOperation) Eval() (float64, error) {
	return m.m.Eval()
}