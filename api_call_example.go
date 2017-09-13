package main

import (
	"fmt"
	"log"

	"github.com/parameterResolver/resolver"
)

func UsageForExtractParametersFromTextApi(service resolver.ISsmParameterService) {
	fmt.Println("Example of ExtractParametersFromText API usage")

	inputDoc := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}"
	resolvedParameters, err := resolver.ExtractParametersFromText(service, inputDoc, resolver.ResolveOptions{
		IgnoreSecureParameters: false,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	for ref, param := range resolvedParameters {
		fmt.Printf("Parameter reference %s -> %s\n", ref, param)
	}
	fmt.Println()
}

func UsageForResolveParameterReferenceList(service resolver.ISsmParameterService) {
	fmt.Println("Example of ResolveParameterReferenceList API usage")

	parameterReferences := []string{
		"ssm:/a/b/c/param1",
		"ssm-secure:param2",
	}

	resolvedParameters, err := resolver.ResolveParameterReferenceList(service, parameterReferences, resolver.ResolveOptions{
		IgnoreSecureParameters: false,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	for ref, param := range resolvedParameters {
		fmt.Printf("Parameter reference %s -> %s\n\n", ref, param)
	}
}

func UsageForResolveParametersInText(service resolver.ISsmParameterService) {
	fmt.Println("Example of ResolveParametersInText API usage")

	unresolvedText := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}"
	resolvedText, err := resolver.ResolveParametersInText(service, unresolvedText, resolver.ResolveOptions{
		IgnoreSecureParameters: false,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Unresolved doc: %s\n", unresolvedText)
	fmt.Printf("Resolved doc:   %s\n\n", resolvedText)
}

func UsageForResolveParametersInFile(service resolver.ISsmParameterService) {
	fmt.Println("Example of ResolveParametersInFile API usage")

	inputFilename := "./test-files/test.json"
	outputFilenameWithSecureStrings := "./resolved_secure_output.json"
	err1 := resolver.ResolveParametersInFile(service, inputFilename, outputFilenameWithSecureStrings, resolver.ResolveOptions{
		IgnoreSecureParameters: false,
	})
	if err1 != nil {
		log.Fatal(err1)
		return
	}

	fmt.Println("Check content of the output file with secure parameters resolved: " + outputFilenameWithSecureStrings)

	outputFilenameIgnoredSecureStrings := "./resolved_no_secure_output.json"
	err2 := resolver.ResolveParametersInFile(service, inputFilename, outputFilenameIgnoredSecureStrings, resolver.ResolveOptions{
		IgnoreSecureParameters: true,
	})
	if err2 != nil {
		log.Fatal(err2)
		return
	}

	fmt.Println("Check content of the output file with ignored secure parameters: " + outputFilenameIgnoredSecureStrings)
}

//
// Preconditions: the following two parameters are provisioned in your AWS account
// 		/a/b/c/param1 is of String type
//      param2 is of SecureString type
//
// Also, run aws configure and supply key, secret and AWS region where the parameters
// were created.
//
func main() {

	service, err := resolver.NewService()
	if err != nil {
		log.Fatal(err)
		return
	}

	UsageForExtractParametersFromTextApi(service)

	UsageForResolveParameterReferenceList(service)

	UsageForResolveParametersInText(service)

	UsageForResolveParametersInFile(service)
}
