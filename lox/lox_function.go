package lox

type LoxFunction struct {
	decl *FunctionStmt
	closure *Environment
}

// call implements Callable
func (f *LoxFunction) call(i *Interpreter, args []any) (retval any) {
	funcEnv := NewEnclosedEnv(f.closure)
	for i := 0; i < len(f.decl.params); i++ {
		funcEnv.Define(f.decl.params[i].Lexeme, args[i])
	}

	defer func() {
		if r := recover(); r != nil {
			// Recovering from return statement
			// Exploit named return value
			if ret, ok := r.(Return); ok {
				retval = ret.value
			} else {
				panic(r)
			}
		}

	}()
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
