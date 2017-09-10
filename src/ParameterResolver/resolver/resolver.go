package resolver

import (
	"log"
	"parameterResolver/ssm"
	"parameterResolver/common"
	"strings"
)


func ResolveParametersInFile(settings *common.Settings, service ssm.ISsmParameterService) error {

	errorInInputParams := common.ValidateInput(settings)
	if errorInInputParams != nil {
		log.Fatal(errorInInputParams)
		return errorInInputParams
	}

	errorInFileOrSize := common.ValidateFileAndSize(settings.InputFile)
	if errorInFileOrSize != nil {
		log.Fatal(errorInFileOrSize)
		return errorInFileOrSize
	}

	unresolvedText, err := common.ReadTextFromFile(settings.InputFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	parametersToFetch := parseParametersFromTextIntoMap(unresolvedText)

	parametersWithValues, err := ssm.GetParametersFromSsmParameterStore(service, parametersToFetch, settings.Options)
	if err != nil {
		log.Fatal(err)
		return err
	}

	resolvedText := resolveDocumentWithParamValues(parametersWithValues, unresolvedText)
	err = common.WriteToFile(resolvedText, settings.OutputFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// parses input text by searching for {{ssm:word}} phrases and stores them into a map
// with default values {{ssm:word}}
func parseParametersFromTextIntoMap(text string) map[string]string {
	matchedPhrases := common.ParameterPlaceholder.FindAllStringSubmatch(text, -1)
	matchedSecurePhrases := common.SecureParameterPlaceholder.FindAllStringSubmatch(text, -1)

	parameterNamesDeduped := make(map[string]string)
	for i := 0; i < len(matchedPhrases); i++ {
		parameterNamesDeduped[matchedPhrases[i][1]] = "{{" + common.SSMNonSecurePrefix + matchedPhrases[i][1] + "}}"
	}

	for i := 0; i < len(matchedSecurePhrases); i++ {
		parameterNamesDeduped[matchedSecurePhrases[i][1]] = "{{" + common.SSMSecurePrefix + matchedSecurePhrases[i][1] + "}}"
	}

	return parameterNamesDeduped
}

func resolveDocumentWithParamValues(parameters map[string]string, text string) string {
	matchedPhrases := common.ParameterPlaceholder.FindAllString(text, -1)
	for i := 0; i < len(matchedPhrases); i++ {
		parameterName := cleanUpParameterPlaceholder(matchedPhrases[i])
		text = strings.Replace(text, matchedPhrases[i], parameters[parameterName], -1)
	}

	matchedSecurePhrases := common.SecureParameterPlaceholder.FindAllString(text, -1)
	for i := 0; i < len(matchedSecurePhrases); i++ {
		parameterName := cleanUpParameterPlaceholder(matchedSecurePhrases[i])
		text = strings.Replace(text, matchedSecurePhrases[i], parameters[parameterName], -1)
	}

	return text
}

func cleanUpParameterPlaceholder(parameterPlaceholder string) string {
	parameterPlaceholder = strings.Join(strings.Fields(parameterPlaceholder),"")

	if strings.HasPrefix(parameterPlaceholder, "{{" + common.SSMNonSecurePrefix) {
		parameterPlaceholder = parameterPlaceholder[len(common.SSMNonSecurePrefix) + 2: len(parameterPlaceholder)-2]
	} else if strings.HasPrefix(parameterPlaceholder, "{{" + common.SSMSecurePrefix) {
		parameterPlaceholder = parameterPlaceholder[len(common.SSMSecurePrefix) + 2: len(parameterPlaceholder)-2]
	}

	return parameterPlaceholder
}
