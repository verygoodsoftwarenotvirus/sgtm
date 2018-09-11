package speakers

type Speaker interface {
	GenerateSpeech(text, fileName string) error
}
