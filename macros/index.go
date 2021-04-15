package main

import (
	"github.com/maxmoehl/calc/types"
)

// Index maps the identifier of plugins to the name of their initialization method
var Index types.Index = map[string]string{
	"sqrt": "NewSqrt",
	"pow": "NewPow",
}
