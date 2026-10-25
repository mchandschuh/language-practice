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

func parse(expression string) (MathExpression, error) {
	return (&Parser{expression: expression}).Parse()
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
		left = BinaryExpression{Operator: operator, Left: left, Right: right}
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
		left = BinaryExpression{Operator: operator, Left: left, Right: right}
	}
}

func (p *Parser) parseUnary() (MathExpression, error) {
	p.skipWhitespace()

	switch p.peek() {
	case '+':
		p.consume()
		return p.parseUnary()
	case '-':
		p.consume()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return UnaryExpression{Operator: Negate, Operand: operand}, nil
	default:
		return p.parsePrimary()
	}
}

func (p *Parser) parsePrimary() (MathExpression, error) {
	p.skipWhitespace()
	if p.eof() {
		return nil, p.errorf("expected a number or parenthesized expression")
	}

	if p.peek() == '(' {
		return p.parseParenthesis()
	}
	if !isDigit(p.peek()) && p.peek() != '.' {
		return nil, p.errorf("expected a number or parenthesized expression; got %q", p.peek())
	}

	value, err := p.parseNumber()
	if err != nil {
		return nil, err
	}
	return &NumberExpression{Value: value}, nil
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

func (p *Parser) consume() byte {
	if p.eof() {
		return 0
	}
	c := p.expression[p.position]
	p.position++
	return c
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
