package calc

import (
	"fmt"
	"math"
)

func init() {
	err := RegisterMacro("pow", NewPow)
	if err != nil {
		panic(err.Error())
	}
}

type Pow struct {
	base, exp Operation
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

func NewPow(parameter ...Operation) (Macro, error) {
	if len(parameter) != 2 {
		return nil, fmt.Errorf("expected two arguments but got %d argument(s)", len(parameter))
	}
	return &Pow{
		base: parameter[0],
		exp: parameter[1],
	}, nil
}