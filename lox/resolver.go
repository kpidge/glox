package lox

type Resolver struct {
	interpreter *Interpreter
	scopes      *Stack[map[string]bool]
}

func NewResolver(i *Interpreter) *Resolver {
	return &Resolver{interpreter: i, scopes: NewStack[map[string]bool]()}
}

// visitAssignExpr implements ExprVisitor.
func (*Resolver) visitAssignExpr(*Assign) {
	panic("unimplemented")
}

// visitBinaryExpr implements ExprVisitor.
func (*Resolver) visitBinaryExpr(*Binary) {
	panic("unimplemented")
}

// visitCallExpr implements ExprVisitor.
func (*Resolver) visitCallExpr(*Call) {
	panic("unimplemented")
}

// visitGroupingExpr implements ExprVisitor.
func (*Resolver) visitGroupingExpr(*Grouping) {
	panic("unimplemented")
}

// visitLiteralExpr implements ExprVisitor.
func (*Resolver) visitLiteralExpr(*Literal) {
	panic("unimplemented")
}

// visitLogicalExpr implements ExprVisitor.
func (*Resolver) visitLogicalExpr(*Logical) {
	panic("unimplemented")
}

// visitUnaryExpr implements ExprVisitor.
func (*Resolver) visitUnaryExpr(*Unary) {
	panic("unimplemented")
}

// visitVariableExpr implements ExprVisitor.
func (r *Resolver) visitVariableExpr(expr *Variable) {
	scope, err := r.scopes.Peek()
	if err == nil && !scope[expr.Name.Lexeme] {
		ErrorOnToken(expr.Name, "Can't read local variable in its own initialiser")
	}
	r.resolveLocal(expr, expr.Name)
}

// visitBlockStmt implements StmtVisitor.
func (r *Resolver) visitBlockStmt(stmt *BlockStmt) {
	r.beginScope()
	r.resolve(stmt.statements)
	r.endScope()
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name Token) {
	// Are we getting a reference to the map[string]bool on the stack,
	// or a copy?
	scope, err := r.scopes.Peek()
	if err != nil {
		return
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	scope, err := r.scopes.Peek()
	if err != nil {
		return
	}
	// Variable is initialised
	scope[name.Lexeme] = true
}

func (r *Resolver) resolve(statements []Stmt) {
	for _, stmt := range statements {
		resolveStmt(stmt)
	}
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	// TODO: Walk stack to find variable and then resolve it
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

// visitExpressionStmt implements StmtVisitor.
func (*Resolver) visitExpressionStmt(*ExpressionStmt) {
	panic("unimplemented")
}

// visitFunctionStmt implements StmtVisitor.
func (*Resolver) visitFunctionStmt(*FunctionStmt) {
	panic("unimplemented")
}

// visitIfStmt implements StmtVisitor.
func (*Resolver) visitIfStmt(*IfStmt) {
	panic("unimplemented")
}

// visitPrintStmt implements StmtVisitor.
func (*Resolver) visitPrintStmt(*PrintStmt) {
	panic("unimplemented")
}

// visitReturnStmt implements StmtVisitor.
func (*Resolver) visitReturnStmt(*ReturnStmt) {
	panic("unimplemented")
}

// visitVarStmt implements StmtVisitor.
func (r *Resolver) visitVarStmt(stmt *VarStmt) {
	r.declare(stmt.Name)
	if init := stmt.Initialiser; init != nil {
		r.resolveExpr(init)
	}
	r.define(stmt.Name)
}

// visitWhileStmt implements StmtVisitor.
func (*Resolver) visitWhileStmt(*WhileStmt) {
	panic("unimplemented")
}
