package lox


type Expr interface {
	Accept(ExprVisitor)
}

type ExprVisitor interface {
	visitBinaryExpr(*Binary)
	visitUnaryExpr(*Unary)
	visitLiteralExpr(*Literal)
	visitGroupingExpr(*Grouping)
	visitVariableExpr(*Variable)
	visitAssignExpr(*Assign)
	visitLogicalExpr(*Logical)
	visitCallExpr(*Call)
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

type Variable struct {
	Name Token
}

func (va *Variable) Accept(v ExprVisitor) {
	v.visitVariableExpr(va)
}

type Assign struct {
	Name Token
	Value Expr
}

func (a *Assign) Accept(v ExprVisitor) {
	v.visitAssignExpr(a)
}

type Logical struct {
	op Token
	left Expr
	right Expr
}

func (l *Logical) Accept(v ExprVisitor) {
	v.visitLogicalExpr(l)
}

type Call struct {
	paren Token
	callee Expr
	arguments []Expr
}

func (c *Call) Accept(v ExprVisitor) {
	v.visitCalleeExpr(c)
}

type Stmt interface {
	Accept(StmtVisitor)
}

type StmtVisitor interface {
	// Expression statement
	visitExpressionStmt(*ExpressionStmt)

	// Print statement
	visitPrintStmt(*PrintStmt)

	// Variable declaration statement
	visitVarStmt(*VarStmt)

	// Block statement
	visitBlockStmt(*BlockStmt)

	// If statement
	visitIfStmt(*IfStmt)

	// While statement
	visitWhileStmt(*WhileStmt)
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

type VarStmt struct {
	Name Token
	Initialiser Expr
}

func (vs *VarStmt) Accept(v StmtVisitor) {
	v.visitVarStmt(vs)
}

type BlockStmt struct {
	statements []Stmt
}

func (bs *BlockStmt) Accept(v StmtVisitor) {
	v.visitBlockStmt(bs)
}

type IfStmt struct {
	expr Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (is *IfStmt) Accept(v StmtVisitor) {
	v.visitIfStmt(is)
}

type WhileStmt struct {
	expr Expr
	body Stmt
}

func (ws *WhileStmt) Accept(v StmtVisitor) {
	v.visitWhileStmt(ws)
}
