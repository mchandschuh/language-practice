package calc

import (
	"fmt"
	"strconv"
)

type Parser struct {
	expression string
	position   int
}

func (p *Parser) Parse() (MathExpression, error) {
	p.skipWhitespace()
	if p.eof() {
		return nil, p.errorf("empty expression")
	}

	result, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipWhitespace()
	if !p.eof() {
		return nil, p.errorf("unexpected character %q", p.peek())
	}

	return result, nil
}

// Parse parses an arithmetic expression into its expression tree.
func Parse(expression string) (MathExpression, error) {
	return (&Parser{expression: expression}).Parse()
}

func parse(expression string) (MathExpression, error) {
	return Parse(expression)
}

func (p *Parser) parseExpression() (MathExpression, error) {
	return p.parseAdditive()
}

func (p *Parser) parseAdditive() (MathExpression, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}

	for {
		p.skipWhitespace()
		if p.eof() || p.peek() == ')' {
			return left, nil
		}

		var operator BinaryOperator
		switch p.peek() {
		case '+':
			operator = Add
		case '-':
			operator = Subtract
		default:
			return nil, p.errorf("expected +, -, or end of expression; got %q", p.peek())
		}
		p.consume()

		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Operator: operator, Left: left, Right: right}
	}
}

func (p *Parser) parseMultiplicative() (MathExpression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for {
		p.skipWhitespace()
		if p.eof() || p.peek() == ')' || p.peek() == '+' || p.peek() == '-' {
			return left, nil
		}

		var operator BinaryOperator
		switch p.peek() {
		case '*':
			operator = Multiply
		case '/':
			operator = Divide
		default:
			return nil, p.errorf("expected *, /, or lower-precedence operator; got %q", p.peek())
		}
		p.consume()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Operator: operator, Left: left, Right: right}
	}
}

func (p *Parser) parseUnary() (MathExpression, error) {
	p.skipWhitespace()

	op := Positive
	switch {
	case p.tryConsumeIf('+'):
		if p.tryConsumeIf('+') {
			op = PrefixIncrement
		}
	case p.tryConsumeIf('-'):
		op = Negate
		if p.tryConsumeIf('-') {
			op = PrefixDecrement
		}
	default:
		return p.parsePostfix()
	}

	operand, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	return &UnaryExpression{Operator: op, Operand: operand}, nil
}

func (p *Parser) parsePostfix() (MathExpression, error) {
	operand, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if p.tryConsumePair('+', '+') {
		return &UnaryExpression{Operator: PostfixIncrement, Operand: operand}, nil
	}
	if p.tryConsumePair('-', '-') {
		return &UnaryExpression{Operator: PostfixDecrement, Operand: operand}, nil
	}
	return operand, nil
}

func (p *Parser) parsePrimary() (MathExpression, error) {
	p.skipWhitespace()
	if p.eof() {
		return nil, p.errorf("expected a number or parenthesized expression")
	}

	if p.peek() == '(' {
		return p.parseParenthesis()
	}
	if isIdentifierStart(p.peek()) {
		return p.parseIdentifier()
	}
	if !isDigit(p.peek()) && p.peek() != '.' {
		return nil, p.errorf("expected an identifier, number, or parenthesized expression; got %q", p.peek())
	}

	value, err := p.parseNumber()
	if err != nil {
		return nil, err
	}
	return &NumberExpression{Value: value}, nil
}

func (p *Parser) parseIdentifier() (MathExpression, error) {
	start := p.position
	p.consume()
	for isIdentifierPart(p.peek()) {
		p.consume()
	}
	return &IdentifierExpression{Name: p.expression[start:p.position]}, nil
}

func (p *Parser) parseParenthesis() (MathExpression, error) {
	p.consume()

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipWhitespace()
	if p.eof() || p.peek() != ')' {
		return nil, p.errorf("expected closing parenthesis")
	}
	p.consume()
	return expression, nil
}

func (p *Parser) parseNumber() (float64, error) {
	start := p.position
	hasDigits := false

	for isDigit(p.peek()) {
		hasDigits = true
		p.consume()
	}

	if p.peek() == '.' {
		p.consume()
		for isDigit(p.peek()) {
			hasDigits = true
			p.consume()
		}
	}

	if !hasDigits {
		return 0, p.errorf("expected digits")
	}

	if p.peek() == 'e' || p.peek() == 'E' {
		p.consume()
		if p.peek() == '+' || p.peek() == '-' {
			p.consume()
		}

		exponentStart := p.position
		for isDigit(p.peek()) {
			p.consume()
		}
		if p.position == exponentStart {
			return 0, p.errorf("expected exponent digits")
		}
	}

	value, err := strconv.ParseFloat(p.expression[start:p.position], 64)
	if err != nil {
		return 0, p.errorf("invalid number %q: %v", p.expression[start:p.position], err)
	}
	return value, nil
}

func (p *Parser) peek() byte {
	if p.eof() {
		return 0
	}
	return p.expression[p.position]
}

func (p *Parser) peekAhead(n int) byte {
	if p.eof() {
		return 0
	}
	if p.position+n >= len(p.expression) {
		return 0
	}
	return p.expression[p.position+n]
}

func (p *Parser) consume() byte {
	if p.eof() {
		return 0
	}
	c := p.expression[p.position]
	p.position++
	return c
}

func (p *Parser) tryConsumeIf(c byte) bool {
	if p.peek() != c {
		return false
	}
	p.position++
	return true
}

func (p *Parser) tryConsumePair(first, second byte) bool {
	if p.peek() != first {
		return false
	}
	if p.peekAhead(1) != second {
		return false
	}
	p.position += 2
	return true
}

func (p *Parser) skipWhitespace() {
	for p.position < len(p.expression) {
		switch p.expression[p.position] {
		case ' ', '\t', '\n', '\r':
			p.position++
		default:
			return
		}
	}
}

func (p *Parser) eof() bool {
	return p.position >= len(p.expression)
}

func (p *Parser) errorf(format string, args ...any) error {
	return fmt.Errorf("at position %d: %s", p.position, fmt.Sprintf(format, args...))
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isIdentifierStart(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isIdentifierPart(c byte) bool {
	return isIdentifierStart(c) || isDigit(c)
}
