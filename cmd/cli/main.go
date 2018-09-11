package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/interpreter"
	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/speakers/say"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

var (
	cfgFile         string
	verbose         bool
	filePath        string
	functionsToRead []string
	speakerToRead	string
	parts           map[string]struct{}

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
		Long:  "reads a file of your choice",
		RunE: func(*cobra.Command, []string) error {
			x := interpret.NewInterpreter(functionsToRead)
			if err := x.InterpretFile(parseCode(), functionsToRead); err != nil {
				log.Fatal(err)
			}

			y := x.RawOutput()
			fmt.Println(y)
			speaker := say.New("en", speakerToRead)
			speaker.GenerateSpeech(y, "")
			return nil
		},
	}
	readCommand.Flags().StringVarP(&filePath, "file", "f", "", "the file you want to read")
	readCommand.Flags().StringArrayVarP(&functionsToRead, "function", "p", nil, "the functions you want to read from the file")
	readCommand.Flags().StringVarP(&speakerToRead, "speaker", "s", "Alex", "the speaker you want to read")


	rootCmd.AddCommand(readCommand)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
