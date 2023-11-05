package ast

import (
	"bytes"

	"go_interp/model/token"
)

type Node interface {
	// TokenLiteral returns the literal string which was parsed to create this node.
	TokenLiteral() string
	String() string
}

// Statement no result value
type Statement interface {
	Node
	statementNode()
}

// Expression has result value
type Expression interface {
	Node
	expressionNode()
}

type IntegerLiteral struct {
	Token token.Token // token.INT 词法单元
	Value int64
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token // 前缀词法单元, 如!, -等
	Operator string
	Right    Expression
}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) expressionNode() {}

func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // 运算词法单元, 如+, -等
	Operator string
	Left     Expression
	Right    Expression
}

func (p *InfixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *InfixExpression) expressionNode() {}

func (p *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Left.String())
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}
