package common

import (
	"flag"
	"os"
	"strings"
	"errors"
)

type Settings struct  {
	InputFile string
	OutputFile string
	OutputFormat string
	Options string
}

func ParseCommandLine() (s Settings, err error) {

	var inputFile = flag.String("in", "", "Input file")
	var outputFile = flag.String("out", "", "Resolved file")
	var outputFormat = flag.String("fmt", TextOutputFormat, "Output encoding. txt (default), xml, yml, json.")
	var options = flag.String("opt", FailOnParameterNotFoundOption, "Enter failonparameternotfound (default) or ignoreparameternotfound")

	flag.Parse()

	s.InputFile = *inputFile
	s.OutputFile = *outputFile
	s.OutputFormat = *outputFormat
	s.Options = *options

	err = ValidateInput(&s)

	return
}

func ValidateInput(settings *Settings) error {

	if len(settings.InputFile) == 0 {
		return errors.New("Input file name is not provided.")
	}

	if len(settings.OutputFile) == 0 {
		return errors.New("Output file name is not provided.")
	}

	// check if the input file is valid
	_, err := os.Open(settings.InputFile)
	if err != nil {
		return err
	}

	// check if output document format is valid
	_, formatMap := formatMap[strings.ToLower(settings.OutputFormat)]
	if !formatMap {
		return errors.New("Invalid output document format is provided.")
	}

	// check if option is valid
	_, optionMap := optionMap[strings.ToLower(settings.Options)]
	if !optionMap {
		return errors.New("Invalid document resolution option provided.")
	}

	return nil
}