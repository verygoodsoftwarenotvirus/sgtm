package interpret

type Describer interface {
	Describe() (string, error)
}
