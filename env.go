package main

type Environment struct {
	vars   map[Name]Expression
	parent *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		vars:   make(map[Name]Expression),
		parent: parent,
	}
}

func (env *Environment) Get(name Name) (Expression, bool) {
	if value, ok := env.vars[name]; ok {
		return value, true
	}
	if env.parent != nil {
		return env.parent.Get(name)
	}
	return nil, false
}

func (env *Environment) Set(name Name, value Expression) {
	env.vars[name] = value
}
