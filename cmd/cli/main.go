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
	verbose         bool
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
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Help message for toggle")

	readCommand := &cobra.Command{
		Use:   "read",
		Short: "",
		Long:  "reads a file (or parts of a file) of your choice",
		RunE: func(*cobra.Command, []string) error {
			x := interpret.NewInterpreter(functionsToRead)
			if err := x.InterpretFile(parseCode(), functionsToRead); err != nil {
				log.Fatal(err)
			}

			y := x.RawOutput()
			fmt.Println(y)

			switch strings.ToLower(voiceService) {
			case amazonVoiceService:
				speaker = polly.New(speakerToRead)
				return speaker.GenerateSpeech(y, "")
			default:
				speaker = say.New(say.DefaultLanguage, speakerToRead)
				return speaker.GenerateSpeech(y, "")
			}
		},
	}
	readCommand.Flags().StringVarP(&filePath, "file", "f", "", "the file you want to read")
	readCommand.Flags().StringArrayVarP(&functionsToRead, "part", "p", nil, "the functions you want to read from the file")
	readCommand.Flags().StringVarP(&voiceService, "voice", "v", defaultVoiceService, "the TTS service you want to use")
	readCommand.Flags().StringVarP(&speakerToRead, "speaker", "s", say.DefaultVoice, "the speaker you want to read")

	rootCmd.AddCommand(readCommand)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
