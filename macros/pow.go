package main

import (
	"fmt"
	"math"

	"github.com/maxmoehl/calc/types"
)

type Pow struct {
	base, exp types.Node
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

func newPow(parameters []types.Node) (types.Macro, error) {
	if len(parameters) != 2 {
		return nil, fmt.Errorf("expected two arguments but got %d argument(s)", len(parameters))
	}
	return &Pow{
		base: parameters[0],
		exp:  parameters[1],
	}, nil
}
