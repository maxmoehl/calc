package main

import (
	"fmt"
	"math"

	"github.com/maxmoehl/calc/types"
)

type Sqrt struct {
	value types.Node
}

func (s *Sqrt) Eval() (float64, error) {
	f, err := s.value.Eval()
	if err != nil {
		return math.NaN(), err
	}
	return math.Sqrt(f), nil
}

var NewSqrt = types.NewMacro(newSqrt)

func newSqrt(parameters []types.Node) (types.Macro, error) {
	if len(parameters) != 1 {
		return nil, fmt.Errorf("expected one argument but got %d arguments", len(parameters))
	}
	return &Sqrt{
		value: parameters[0],
	}, nil
}
