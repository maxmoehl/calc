package calc

type literal struct {
	value float64
}

func (l *literal) Locked() bool {
	return true
}

func (l *literal) Eval() (float64, error) {
	return l.value, nil
}
