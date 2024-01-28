package lox

import (
	"fmt"
)

type ASTPrinter struct{
	result string
}

func (p *ASTPrinter) PrintAST(e Expr) {
	e.Accept(p)
	fmt.Println(p.result)
}

func (p *ASTPrinter) visitBinaryExpr(b *Binary) {
	p.result = p.parenthesise(b.Op.Lexeme, b.Left, b.Right)
}

func (p *ASTPrinter) visitUnaryExpr(u *Unary) {
	p.result = p.parenthesise(u.Op.Lexeme, u.Right)
}

func (p *ASTPrinter) visitLiteralExpr(l *Literal) {
	if l.Value == nil {
		p.result = "nil"
	} else {
		p.result = fmt.Sprint(l.Value)
	}
}

func (p *ASTPrinter) visitGroupingExpr(g *Grouping) {
	p.result = p.parenthesise("group", g.Expr)
}

func (p *ASTPrinter) visitVariableExpr(v *Variable) {
	p.result = p.parenthesise("var", v)
}

func (p *ASTPrinter) visitAssignExpr(a *Assign) {
	p.result = p.parenthesise("assign", a)
}

func (p *ASTPrinter) visitLogicalExpr(l *Logical) {
	p.result = p.parenthesise(l.op.Lexeme, l.left, l.right)
}

func (p *ASTPrinter) visitCallExpr(l *Call) {
	l.callee.Accept(p)
	p.result = p.parenthesise(p.result, l.arguments...)
}

func (p *ASTPrinter) parenthesise(name string, exprs ...Expr) string {
	res := "(" + name
	for _, e := range exprs {
		res += " "
		e.Accept(p)
		res += p.result
	}
	res += ")"
	return res
}
