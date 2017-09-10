package main
//
//import (
//	"log"
//	"time"
//	"parameterResolver/resolver"
//	"flag"
//)

type Settings struct  {
	InputFile string
	OutputFile string
	OutputFormat string
	Options string
}

func main() {
/*
	settings, err := ParseCommandLine()
	if err != nil {
		log.Fatal(err)
		return
	}

	start := time.Now()

	ssmService, err := resolver.NewService()
	if err != nil {
		log.Fatal(err)
		return
	}

	resolver.ResolveParametersInFile(&settings, ssmService)

	elapsed := time.Since(start)
	log.Printf("It took %s to execute.\n", elapsed)*/
}

/*

func ParseCommandLine() (s Settings, err error) {

	var inputFile = flag.String("in", "", "Input file")
	var outputFile = flag.String("out", "", "Resolved file")
	var outputFormat = flag.String("fmt", TextOutputFormat, "Output encoding. txt (default), xml, yml, json.")
	var options = flag.String("opt", FailOnParameterNotFoundOption, "Enter failonparameternotfound or ignoreparameternotfound")

	flag.Parse()

	s.InputFile = *inputFile
	s.OutputFile = *outputFile
	s.OutputFormat = *outputFormat
	s.Options = *options

	err = ValidateInput(&s)

	return
}*/