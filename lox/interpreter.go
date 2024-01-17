package lox

import "fmt"

type Interpreter struct {
	tmp any
}

type RuntimeError struct {
	token Token
	msg string
}

func (err RuntimeError) Error() string {
	return err.msg
}

func (i *Interpreter) Interpret(expr Expr) {
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
	value := i.evaluate(expr)
	fmt.Println(stringify(value))
}


func (i *Interpreter) evaluate(expr Expr) any {
	expr.Accept(i)
	return i.tmp
}
	
func (i *Interpreter) visitBinaryExpr(expr *Binary) {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)
	fmt.Println(left)
	fmt.Println(right)

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

func (i *Interpreter) isTruthy(obj any) bool {
	if obj == nil { return false }
	
	if v, ok := obj.(bool); ok {
		return v
	}
	return true
}

func (i *Interpreter) isEqual(left any, right any) bool {
	if left == nil && right == nil { return true }
	if left == nil { return false }

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
