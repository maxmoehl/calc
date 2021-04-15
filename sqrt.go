package calc

import (
	"fmt"
	"math"
)

func init() {
	err := RegisterMacro("sqrt", NewSqrt)
	if err != nil {
		panic(err.Error())
	}
}

type Sqrt struct {
	value Operation
}

func (s *Sqrt) Eval() (float64, error) {
	f, err := s.value.Eval()
	if err != nil {
		return math.NaN(), err
	}
	return math.Sqrt(f), nil
}

func NewSqrt(parameter ...Operation) (Macro, error) {
	if len(parameter) != 1 {
		return nil, fmt.Errorf("expected one argument but got %d arguments", len(parameter))
	}
	return &Sqrt{
		value: parameter[0],
	}, nil
}