package lox

type Environment struct {
	enclosing *Environment
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]any)}
}

func NewEnclosedEnv(encl *Environment) *Environment {
	return &Environment{enclosing: encl, values : make(map[string]any)}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) any {
	if value, ok := e.values[name.Lexeme]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		} else {
			panic(RuntimeError{token: name, msg: "Undefined variable '" + name.Lexeme + "'."})
		}
	} else {
		return value
	}
}

func (e *Environment) Assign(name Token, value any) {
	if _, ok := e.values[name.Lexeme]; !ok {
		if e.enclosing != nil {
			e.enclosing.Assign(name, value)
		} else {
			panic(RuntimeError{token: name, msg: "Undefined variable '" + name.Lexeme + "'."})
		}
	}
	e.Define(name.Lexeme, value)
}
