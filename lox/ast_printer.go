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
