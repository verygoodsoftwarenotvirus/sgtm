package polly

import (
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/stretchr/testify/assert"
)

const exampleMP3Path = "pkg/speakers/polly.go/test_files/test.mp3"

func loadExampleFile(t *testing.T) io.ReadCloser {
	t.Helper()
	f, err := os.Open(exampleMP3Path)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestPollySpeaker_saveAsMP3(T *testing.T) {
	T.Parallel()

	ps := New("joanna")

	T.Run("generic", func(t *testing.T) {
		t.Parallel()

		exampleOutputFilename := "testsave.mp3"
		exampleData := &polly.SynthesizeSpeechOutput{AudioStream: loadExampleFile(t)}

		err := ps.saveAsMP3(exampleOutputFilename, exampleData)
		assert.NoError(t, err)
		assert.NoError(t, os.Remove(exampleOutputFilename))
	})
}
