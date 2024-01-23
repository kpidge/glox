package lox


type Expr interface {
	Accept(ExprVisitor)
}

type ExprVisitor interface {
	visitBinaryExpr(*Binary)
	visitUnaryExpr(*Unary)
	visitLiteralExpr(*Literal)
	visitGroupingExpr(*Grouping)
}

type Binary struct {
	Left Expr
	Right Expr
	Op Token
}

// Pointer receiver means that types like *Binary implement 
// the Expr interface - do we need to pass around pointers everywhere?
// Currently thinking 'yes'; otherwise, any code that traverses the AST 
// would effectively end up copying the whole tree
func (b *Binary) Accept(v ExprVisitor) {
	v.visitBinaryExpr(b)
}

type Unary struct {
	Op Token
	Right Expr
}

func (u *Unary) Accept(v ExprVisitor) {
	v.visitUnaryExpr(u)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(v ExprVisitor) {
	v.visitLiteralExpr(l)
}

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(v ExprVisitor) {
	v.visitGroupingExpr(g)
}

type Stmt interface {
	Accept(StmtVisitor)
}

type StmtVisitor interface {
	// Expression statement
	visitExpressionStmt(*ExpressionStmt)

	// Print statement
	visitPrintStmt(*PrintStmt)
}

type ExpressionStmt struct {
	Expr Expr
}

func (es *ExpressionStmt) Accept(v StmtVisitor) {
	v.visitExpressionStmt(es)
}

type PrintStmt struct {
	Expr Expr
}

func (ps *PrintStmt) Accept(v StmtVisitor) {
	v.visitPrintStmt(ps)
}
