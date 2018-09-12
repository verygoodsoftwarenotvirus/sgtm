package polly

import (
	"io"
	"os"
	"strings"

	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/speakers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

const (
	DefaultVoice = "joanna"
)

var (
	// acceptableVoices is the set of voices currently supported by SGTM. More AWS Polly voices can be found here: https://docs.aws.amazon.com/polly/latest/dg/voicelist.html.
	acceptableVoices = map[string]*string{
		DefaultVoice: aws.String("Joanna"),
		"joey":       aws.String("Joey"),
		"justin":     aws.String("Justin"),
		"matthew":    aws.String("Matthew"),
		"ivy":        aws.String("Ivy"),
		"kendra":     aws.String("Kendra"),
		"kimberly":   aws.String("Kimberly"),
		"salli":      aws.String("Salli"),
	}
)

var _ speakers.Speaker = (*PollySpeaker)(nil)

type PollySpeaker struct {
	*polly.Polly
	VoiceID *string
}

// New creates an AWS Polly session. It takes in a voiceName string to determine which speaker to use.
func New(voiceName string) *PollySpeaker {
	var voice *string
	voice, ok := acceptableVoices[voiceName]
	if !ok {
		voice = acceptableVoices[DefaultVoice]
	}

	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &PollySpeaker{
		Polly:   polly.New(sess),
		VoiceID: voice,
	}
}

// GenerateSpeech takes in a string and outputs Synthesized Speech and possible an error.
func (ps *PollySpeaker) GenerateSpeech(s string, fileName string) error {
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(s), VoiceId: ps.VoiceID}
	output, err := ps.SynthesizeSpeech(input)
	if err != nil {
		return err
	}
	return ps.saveAsMP3(fileName, output)
}

// SaveAsMP3 takes a file name and the synthesized speech from GenerateSpeech and saves a MP3 file of the speech to that location.
func (ps *PollySpeaker) saveAsMP3(fileName string, output *polly.SynthesizeSpeechOutput) error {
	// Save as MP3
	names := strings.Split(fileName, ".")
	name := names[0]
	mp3File := name + ".mp3"

	outFile, err := os.Create(mp3File)
	if err != nil {
		return err
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		return err
	}
	return nil
}
