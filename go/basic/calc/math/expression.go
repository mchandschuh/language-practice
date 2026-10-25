package calc

type MathExpression interface {
	expressionNode()
}

type NumberExpression struct {
	Value float64
}

func (*NumberExpression) expressionNode() {}

type UnaryExpression struct {
	Operator UnaryOperator
	Operand  MathExpression
}

func (UnaryExpression) expressionNode() {}

type BinaryExpression struct {
	Operator BinaryOperator
	Left     MathExpression
	Right    MathExpression
}

func (BinaryExpression) expressionNode() {}

type MathExpressionJob struct {
	Result chan float64
	Left   MathExpression
	Right  MathExpression
}
