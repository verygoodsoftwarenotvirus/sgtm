package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/interpreter"
	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/speakers"
	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/speakers/polly"
	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/speakers/say"

	"github.com/spf13/cobra"
)

const (
	defaultVoiceService = "say"
	amazonVoiceService  = "polly"
)

var (
	cfgFile,
	filePath,
	speakerToRead,
	voiceService string
	speaker         speakers.Speaker
	parts           map[string]struct{}
	verbose, debug  bool
	functionsToRead []string

	defaultSpeakers = map[string]string{
		defaultVoiceService: say.DefaultVoice,
		amazonVoiceService:  polly.DefaultVoice,
	}

	ErrInvalidFilePath     = errors.New("invalid file path")
	ErrInvalidInstructions = errors.New("invalid instructions")
)

func use(...interface{}) {}

func parseCode() *ast.File {
	code, err := openFile()
	if err != nil {
		log.Fatal(err)
	}

	p, err := parser.ParseFile(token.NewFileSet(), filePath, code, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sgtm",
	Short: "",
	Long:  `SGTM is a screen reader meant for reading computer code out to you`,
}

func openFile() (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "do you want overly wordy output or not")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "print instead of actually playing sounds")

	readCommand := &cobra.Command{
		Use:   "read",
		Short: "",
		Long:  "reads a file (or parts of a file) of your choice",
		RunE: func(*cobra.Command, []string) error {
			if filePath == "" {
				return errors.New("no filename provided!")
			}
			x := interpret.NewInterpreter(functionsToRead, verbose)
			if err := x.InterpretFile(parseCode(), functionsToRead); err != nil {
				log.Fatal(err)
			}

			y := x.RawOutput()

			if debug {
				fmt.Println(y)
				return nil
			}

			switch strings.ToLower(voiceService) {
			case amazonVoiceService:
				speaker = polly.New(speakerToRead)
				return speaker.GenerateSpeech(y, filePath)
			default:
				speaker = say.New(say.DefaultLanguage, speakerToRead)
				return speaker.GenerateSpeech(y, filePath)
			}
		},
	}
	readCommand.Flags().StringVarP(&filePath, "file", "f", "", "the file you want to read")
	readCommand.Flags().StringArrayVarP(&functionsToRead, "part", "p", nil, "the functions you want to read from the file")
	readCommand.Flags().StringVarP(&voiceService, "voice-service", "k", say.DefaultVoice, "the TTS service you want to use")
	readCommand.Flags().StringVarP(&speakerToRead, "speaker", "s", defaultVoiceService, "the speaker you want to read")

	rootCmd.AddCommand(readCommand)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
