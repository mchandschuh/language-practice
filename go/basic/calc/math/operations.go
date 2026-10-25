package calc

type BinaryOperator uint8
type UnaryOperator uint8

const (
	Add BinaryOperator = iota
	Subtract
	Multiply
	Divide
)

const (
	Positive UnaryOperator = iota
	Negate
	PrefixDecrement
	PostfixDecrement
	PrefixIncrement
	PostfixIncrement
)
