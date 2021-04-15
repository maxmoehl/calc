package main

import (
	"fmt"
	"math"

	"github.com/maxmoehl/calc/types"
)

type Sqrt struct {
	value types.Operation
}

func (s *Sqrt) Eval() (float64, error) {
	f, err := s.value.Eval()
	if err != nil {
		return math.NaN(), err
	}
	return math.Sqrt(f), nil
}

var NewSqrt = types.NewMacro(newSqrt)

func newSqrt(parameter []types.Operation) (types.Macro, error) {
	if len(parameter) != 1 {
		return nil, fmt.Errorf("expected one argument but got %d arguments", len(parameter))
	}
	return &Sqrt{
		value: parameter[0],
	}, nil
}