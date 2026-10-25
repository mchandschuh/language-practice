package calc

// MathExpression is the marker interface implemented by every calculator AST
// node. It describes expression structure; evaluation is performed by
// Evaluator rather than by the nodes themselves.
type MathExpression interface {
	expressionNode()
}

// NumberExpression represents a numeric literal.
type NumberExpression struct {
	Value float64
}

func (*NumberExpression) expressionNode() {}

// IdentifierExpression represents a named input value.
type IdentifierExpression struct {
	Name string
}

func (*IdentifierExpression) expressionNode() {}

// UnaryExpression represents a unary operation and its operand.
type UnaryExpression struct {
	Operator UnaryOperator
	Operand  MathExpression
}

func (*UnaryExpression) expressionNode() {}

// BinaryExpression represents a binary operation and its operands.
type BinaryExpression struct {
	Operator BinaryOperator
	Left     MathExpression
	Right    MathExpression
}

func (*BinaryExpression) expressionNode() {}
