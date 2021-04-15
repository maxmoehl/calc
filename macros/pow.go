package main

import (
	"fmt"
	"math"

	"github.com/maxmoehl/calc/types"
)

type Pow struct {
	base, exp types.Operation
}

func (p *Pow) Eval() (float64, error) {
	base, err := p.base.Eval()
	if err != nil {
		return math.NaN(), err
	}
	exp, err := p.exp.Eval()
	if err != nil {
		return math.NaN(), err
	}
	return math.Pow(base, exp), nil
}

var NewPow = types.NewMacro(newPow)

func newPow(parameter []types.Operation) (types.Macro, error) {
	if len(parameter) != 2 {
		return nil, fmt.Errorf("expected two arguments but got %d argument(s)", len(parameter))
	}
	return &Pow{
		base: parameter[0],
		exp: parameter[1],
	}, nil
}