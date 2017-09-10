package main

import (
	"fmt"
	"log"
	"parameterResolver/resolver"
)

func UsageForExtractParametersFromTextApi(service resolver.ISsmParameterService) {
	inputDoc := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}"
	resolvedParameters, err := resolver.ExtractParametersFromText(service, inputDoc, resolver.ResolveOptions{
		ResolveSecureParameters:true,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	for ref, param := range resolvedParameters {
		fmt.Printf("-- For reference %s\n %s\n\n", ref, param)
		fmt.Println()
	}
}

func main() {

	service, err := resolver.NewService()
	if err != nil {
		log.Fatal(err)
		return
	}

	UsageForExtractParametersFromTextApi(service)

	fmt.Println()

}
