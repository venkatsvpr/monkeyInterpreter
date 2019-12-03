package ast

import (
	"monkey/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _,s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// let statementNode
type LetStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if (ls != nil) {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// identifer
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal}
func (i *Identifier) String() string {
	return i.Value
}


// integer literal
type IntegerLiteral struct {
	Value int64
	Token token.Token
}
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// return statement
type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {return rs.Token.Literal}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if (rs.ReturnValue != nil) {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// expression Statement
type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {return es.Token.Literal}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// prefix expression
type PrefixExpression struct {
	Token token.Token
	Operator string
	Right Expression
}

func (ps *PrefixExpression) expressionNode() {}
func (ps *PrefixExpression) TokenLiteral() string { return ps.Token.Literal }
func (ps *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ps.Operator)
	out.WriteString(ps.Right.String())
	out.WriteString(")")
	return out.String()
}

// infix expression
type InfixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String()+" "+ie.Operator+" "+ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// boolean
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string { return b.Token.Literal }

