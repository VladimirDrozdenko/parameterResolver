package resolver

import (
	"errors"
	"regexp"
	"strings"
)

//
// Takes text document and resolves all parameters in it according to ResolveOptions.
// It will return a map of (parameter reference) to SsmParameterInfo.
func ExtractParametersFromText(
	service ISsmParameterService,
	input string,
	options ResolveOptions) (map[string]SsmParameterInfo, error) {

	uniqueParameterReferences, err := parseParametersFromTextIntoMap(input)
	if err != nil {
		return nil, err
	}

	parametersWithValues, err := getParametersFromSsmParameterStore(service, uniqueParameterReferences)
	if err != nil {
		return nil, err
	}

	prefixValidationError := validateParameterReferencePrefix(&parametersWithValues)
	if prefixValidationError != nil {
		return nil, prefixValidationError
	}

	if options.IgnoreSecureParameters {
		for key, value := range parametersWithValues {
			if strings.HasPrefix(key, ssmSecurePrefix) || value.Type == secureStringType {
				delete(parametersWithValues, key)
			}
		}
	}

	encodeResolvedValues(&parametersWithValues, options)

	return parametersWithValues, nil
}

//
// Takes a list of references to SSM parameters, resolves them according to ResolveOptions and
// returns a map of (parameter reference) to SsmParameterInfo.
func ResolveParameterReferenceList(
	service ISsmParameterService,
	parameterReferences []string,
	options ResolveOptions) (map[string]SsmParameterInfo, error) {

	uniqueParameterReferences := dedupSlice(parameterReferences)
	parametersWithValues, err := getParametersFromSsmParameterStore(service, uniqueParameterReferences)
	if err != nil {
		return nil, err
	}

	prefixValidationError := validateParameterReferencePrefix(&parametersWithValues)
	if prefixValidationError != nil {
		return nil, prefixValidationError
	}

	if options.IgnoreSecureParameters {
		for key, value := range parametersWithValues {
			if strings.HasPrefix(key, ssmSecurePrefix) || value.Type == secureStringType {
				delete(parametersWithValues, key)
			}
		}
	}

	encodeResolvedValues(&parametersWithValues, options)

	return parametersWithValues, nil
}

//
// Takes text document, resolves all parameters in it according to ResolveOptions
// and returns resolved document.
func ResolveParametersInText(
	service ISsmParameterService,
	input string,
	options ResolveOptions) (string, error) {

	resolvedParametersMap, err := ExtractParametersFromText(service, input, options)
	if err != nil || resolvedParametersMap == nil || len(resolvedParametersMap) == 0 {
		return input, err
	}

	for ref, param := range resolvedParametersMap {
		var placeholder = regexp.MustCompile("{{\\s*" + ref + "\\s*}}")
		input = placeholder.ReplaceAllString(input, param.Value)
	}

	return input, nil
}

//
// Reads inputFileName, resolves SSM parameters in it according to ResolveOptions and
// stores resolved document in the outputFileName file.
func ResolveParametersInFile(
	service ISsmParameterService,
	inputFileName string,
	outputFileName string,
	options ResolveOptions) error {

	if len(inputFileName) == 0 {
		return errors.New("input file name is not provided")
	}

	if len(outputFileName) == 0 {
		return errors.New("output file name is not provided")
	}

	errorInFileOrSize := validateFileAndSize(inputFileName)
	if errorInFileOrSize != nil {
		return errorInFileOrSize
	}

	unresolvedText, err := readTextFromFile(inputFileName)
	if err != nil {
		return err
	}

	resolvedParametersMap, err := ExtractParametersFromText(service, unresolvedText, options)
	if err != nil || resolvedParametersMap == nil || len(resolvedParametersMap) == 0 {
		return err
	}

	for ref, param := range resolvedParametersMap {
		var placeholder = regexp.MustCompile("{{\\s*" + ref + "\\s*}}")
		unresolvedText = placeholder.ReplaceAllString(unresolvedText, param.Value)
	}

	err = writeToFile(unresolvedText, outputFileName)
	if err != nil {
		return err
	}

	return nil
}

func validateParameterReferencePrefix(resolvedParametersMap *map[string]SsmParameterInfo) error {
	for key, value := range *resolvedParametersMap {
		if strings.HasPrefix(key, ssmSecurePrefix) && value.Type != secureStringType {
			return errors.New("secure prefix " + ssmSecurePrefix + " is used for a non-secure type " + value.Type)
		}

		if strings.HasPrefix(key, ssmNonSecurePrefix) && value.Type == secureStringType {
			return errors.New("non-secure prefix " + ssmNonSecurePrefix + " is used for a secure type " + value.Type)
		}
	}

	return nil
}

func encodeResolvedValues(resolvedParametersMap *map[string]SsmParameterInfo, options ResolveOptions) {
	formatEncoder := NewFormatEncoder(options.ValueEncoding)
	for key, _ := range *resolvedParametersMap {
		(*resolvedParametersMap)[key] = SsmParameterInfo{
			Name:  (*resolvedParametersMap)[key].Name,
			Type:  (*resolvedParametersMap)[key].Type,
			Value: formatEncoder.encode((*resolvedParametersMap)[key].Value),
		}
	}
}

func dedupSlice(slice []string) []string {
	ht := map[string]bool{}

	for _, element := range slice {
		ht[element] = true
	}

	keys := make([]string, len(ht))

	i := 0
	for k := range ht {
		keys[i] = k
		i++
	}

	return keys
}

func parseParametersFromTextIntoMap(text string) ([]string, error) {
	matchedPhrases := parameterPlaceholder.FindAllStringSubmatch(text, -1)
	matchedSecurePhrases := secureParameterPlaceholder.FindAllStringSubmatch(text, -1)

	parameterNamesDeduped := make(map[string]bool)
	for i := 0; i < len(matchedPhrases); i++ {
		parameterNamesDeduped[matchedPhrases[i][1]] = true
	}

	for i := 0; i < len(matchedSecurePhrases); i++ {
		parameterNamesDeduped[matchedSecurePhrases[i][1]] = true
	}

	result := []string{}
	for key, _ := range parameterNamesDeduped {
		result = append(result, key)
	}

	return result, nil
}
