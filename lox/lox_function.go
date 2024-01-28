package lox

type LoxFunction struct {
	decl *FunctionStmt
}

// call implements Callable
func (f *LoxFunction) call(i *Interpreter, args []any) any {
	funcEnv := NewEnclosedEnv(i.env)
	for i := 0; i < len(f.decl.params); i++ {
		funcEnv.Define(f.decl.params[i].Lexeme, args[i])
	}
	i.executeBlock(f.decl.body, funcEnv)
	return nil
}

// arity implements Callable
func (f *LoxFunction) arity() int {
	return len(f.decl.params)
}

// String implements Stringer
func (f *LoxFunction) String() string {
	return "<fn " + f.decl.name.Lexeme + ">"
}
