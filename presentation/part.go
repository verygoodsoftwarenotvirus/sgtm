type Interpreter interface {
	InterpretFile(input *ast.File, chunks []string) error
	RawOutput() string
}
