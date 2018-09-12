type Interpreter interface {
	InterpretFile(input *ast.File, chunks []string) error
	RawOutput() string
}

/*
	blah blah blah
*/

type Describer interface {
	Describe() (string, error)
	GetName() string
}
