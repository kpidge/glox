package lox

import (
	"fmt"
)

type Interpreter struct {
	tmp any
	globals Environment
	env *Environment
}

func NewInterpreter() *Interpreter {
	var globals Environment
	// globals.values["clock"] =

	// Could have env enclose the global environment?
	return &Interpreter{globals: globals, env: &globals}
}

type Callable interface {
	call(*Interpreter, []any) any
	arity() int
}

// TODO: Improve errors and error handling
type RuntimeError struct {
	token Token
	msg   string
}

func (err RuntimeError) Error() string {
	return err.msg
}

// Not actually an error - used for breaking out of e.g. functions 
// with a 'return' statement
type Return struct {
	value any
}

func (ret Return) Error() string {
	return stringify(ret.value)
}

func (i *Interpreter) Interpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			// Determine that we are recovering from a
			// runtimeError
			if re, ok := r.(RuntimeError); ok {
				ReportRuntimeError(re)
			} else {
				panic(r)
			}
		}

	}()
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) evaluate(expr Expr) any {
	expr.Accept(i)
	return i.tmp
}

// visitVarStmt implements StmtVisitor.
func (i *Interpreter) visitVarStmt(stmt *VarStmt) {
	var value any
	if init := stmt.Initialiser; init != nil {
		value = i.evaluate(init)
	}
	i.env.Define(stmt.Name.Lexeme, value)
}

// visitAssignExpr implements StmtVisitor.
func (i *Interpreter) visitAssignExpr(expr *Assign) {
	value := i.evaluate(expr.Value)
	i.env.Assign(expr.Name, value)
	i.tmp = value
}

func (i *Interpreter) visitExpressionStmt(stmt *ExpressionStmt) {
	i.evaluate(stmt.Expr)
}

func (i *Interpreter) visitPrintStmt(stmt *PrintStmt) {
	value := i.evaluate(stmt.Expr)
	fmt.Println(stringify(value))
}

func (i *Interpreter) visitBlockStmt(stmt *BlockStmt) {
	blockEnv := NewEnclosedEnv(i.env)
	i.executeBlock(stmt.statements, blockEnv)
}

func (i *Interpreter) visitIfStmt(stmt *IfStmt) {
	if i.isTruthy(i.evaluate(stmt.expr)) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}
}

func (i *Interpreter) visitWhileStmt(stmt *WhileStmt) {
	for i.isTruthy(i.evaluate(stmt.expr)) {
		i.execute(stmt.body)
	}
}

func (i *Interpreter) visitReturnStmt(stmt *ReturnStmt) {
	var value any
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	}
	panic(Return{value: value})
}

func (i *Interpreter) executeBlock(statements []Stmt, env *Environment) {
	prev := i.env
	defer func() {
		// Reset state of interpreter
		i.env = prev
	}()
	i.env = env
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) visitCallExpr(expr *Call) {
	callee := i.evaluate(expr.callee)

	var args []any
	for _, a := range expr.arguments {
		args = append(args, i.evaluate(a))
	}

	if function, ok := callee.(Callable); !ok {
		panic(RuntimeError{token: expr.paren, msg: "Not a callable expression"})
	} else {
		if len(args) != function.arity() {
			msg := fmt.Sprintf("Expected %d argument but received %d", function.arity(), len(args))
			panic(RuntimeError{token: expr.paren, msg: msg})
		}
		i.tmp = function.call(i, args)
	}
}

func (i *Interpreter) visitFunctionStmt(stmt *FunctionStmt) {
	function := &LoxFunction{decl: stmt, closure: i.env}
	i.env.Define(stmt.name.Lexeme, function)
}

func (i *Interpreter) visitBinaryExpr(expr *Binary) {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Op.Type {
	case GREATER:
		// Number literals are parsed as float64s
		// by the Scanner
		checkNumberOperands(expr.Op, left, right)
		i.tmp = left.(float64) > right.(float64)
	case GREATER_EQUAL:
		checkNumberOperands(expr.Op, left, right)
		i.tmp = left.(float64) >= right.(float64)
	case LESS:
		checkNumberOperands(expr.Op, left, right)
		i.tmp = left.(float64) < right.(float64)
	case LESS_EQUAL:
		checkNumberOperands(expr.Op, left, right)
		i.tmp = left.(float64) <= right.(float64)
	case BANG_EQUAL:
		checkNumberOperands(expr.Op, left, right)
		i.tmp = !i.isEqual(left, right)
	case EQUAL_EQUAL:
		checkNumberOperands(expr.Op, left, right)
		i.tmp = i.isEqual(left, right)
	case MINUS:
		checkNumberOperands(expr.Op, left, right)
		i.tmp = left.(float64) - right.(float64)
	case PLUS:
		l, okLeft := left.(float64)
		r, okRight := right.(float64)
		if okLeft && okRight {
			i.tmp = l + r
			break
		}
		lStr, okLeft := left.(string)
		rStr, okRight := right.(string)
		if okLeft && okRight {
			i.tmp = lStr + rStr
			break
		}
		err := RuntimeError{token: expr.Op, msg: "Operands must be two numbers or two strings"}
		panic(err)
	case SLASH:
		checkNumberOperands(expr.Op, left, right)
		i.tmp = left.(float64) / right.(float64)
	case STAR:
		checkNumberOperands(expr.Op, left, right)
		i.tmp = left.(float64) * right.(float64)
	default:
		i.tmp = nil
	}
}

func (i *Interpreter) visitLogicalExpr(expr *Logical) {
	left := i.evaluate(expr.left)
	if expr.op.Type == OR {
		if i.isTruthy(left) {
			i.tmp = left 
			return
		}
	} else {
		if !i.isTruthy(left) {
			i.tmp = left 
			return
		}
	}
	i.tmp = i.evaluate(expr.right)
}

func (i *Interpreter) visitLiteralExpr(expr *Literal) {
	i.tmp = expr.Value
}

func (i *Interpreter) visitGroupingExpr(expr *Grouping) {
	i.tmp = i.evaluate(expr.Expr)
}

func (i *Interpreter) visitUnaryExpr(expr *Unary) {
	i.tmp = i.evaluate(expr.Right)

	switch expr.Op.Type {
	case MINUS:
		checkNumberOperand(expr.Op, i.tmp)
		i.tmp = -i.tmp.(float64)
	case BANG:
		i.tmp = !i.isTruthy(i.tmp)
	}

	i.tmp = nil
}

// visitVariableExpr implements ExprVisitor.
func (i *Interpreter) visitVariableExpr(v *Variable) {
	i.tmp = i.env.Get(v.Name)
}

func (i *Interpreter) isTruthy(obj any) bool {
	if obj == nil {
		return false
	}

	switch obj.(type) {
	case float64:
		return obj != float64(0)
	case string:
		return obj != ""
	case bool:
		return obj.(bool)
	}
	return true
}

func (i *Interpreter) isEqual(left any, right any) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}

	// What is the behaviour of this? Is it different from the
	// Java implementation?
	// FIXME: can probably do away with the separate logic for nil above
	return left == right
}

func checkNumberOperand(token Token, operand any) {
	if _, ok := operand.(float64); !ok {
		err := RuntimeError{token: token, msg: "Operand must be a number"}
		// Panic as we need to unwind call stack
		panic(err)
	}
}

func checkNumberOperands(token Token, left any, right any) {
	_, okLeft := left.(float64)
	_, okRight := right.(float64)
	if !(okLeft && okRight) {
		err := RuntimeError{token: token, msg: "Operands must be numbers"}
		// Panic as we need to unwind call stack
		panic(err)
	}
}

func stringify(obj any) string {
	return fmt.Sprint(obj)
}
