package resolver

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"parameterResolver/ssm"
	"parameterResolver/common"
	"strconv"
	"os"
)

const testOutputFileName = "out.file"

type ServiceAllResolvedMockedObject struct {
	numberOfParameters int
	ssm.ISsmParameterService
}

func NewServiceAllResolvedMockedObject(n int) (mock ServiceAllResolvedMockedObject) {
	return ServiceAllResolvedMockedObject{
		numberOfParameters: n,
	}
}

func (m *ServiceAllResolvedMockedObject) GetParameters(names []string) (map[string]string, []string, error) {
	parameters := map[string]string {}
	for i := 0; i < m.numberOfParameters; i++ {
		name := "name" + strconv.Itoa(i)
		parameters[name] = "value_" + name
	}

	return parameters, nil, nil
}


func TestCleanUpParameterPlaceholder(t *testing.T) {
	assert.True(t, cleanUpParameterPlaceholder("{{ ssm:name 	}}") == "name")
	assert.True(t, cleanUpParameterPlaceholder("{{	ssm:/a/b_c/name}}") == "/a/b_c/name")
}

func TestResolveDocumentWithParamValues(t *testing.T) {
	unresolvedText := `Some text {{ ssm:param1	 }} some more text {{ssm:/a/b_c/name }} and {{ssm:param1}}`
	parameters := map[string]string {
		"param1" : "value1",
		"/a/b_c/name" : "value2",
	}
	expectedResolvedText := "Some text value1 some more text value2 and value1"

	resolvedText := resolveDocumentWithParamValues(parameters, unresolvedText)
	assert.True(t, resolvedText == expectedResolvedText)
}

func TestParseParametersFromTextIntoMap(t *testing.T) {
	text := `Some text {{ ssm:param1	 }} some more text {{ssm:/a/b_c/name }} and {{ssm:param1}} + {{ssm-secure:blah}}`
	expectedMap := map[string]string {
		"param1" : "{{ssm:param1}}",
		"/a/b_c/name" : "{{ssm:/a/b_c/name}}",
		"blah" : "{{ssm-secure:blah}}",
	}

	paramMap := parseParametersFromTextIntoMap(text)

	assert.True(t, reflect.DeepEqual(paramMap, expectedMap))
}


type resolutionTestCase struct {
	testCaseName string
	numberOfParameters int
	setting common.Settings
	resultFileName string
}

func TestResolveParameters(t *testing.T) {

	defer os.Remove(testOutputFileName)

	var resolutionSettings = []resolutionTestCase {
		{
			testCaseName : "Resolve all parameters without paging",
			numberOfParameters : 5,
			setting : common.Settings {
				InputFile: "../test_files/all_resolved_no_paging.json",
				OutputFile: testOutputFileName,
				OutputFormat: common.TextOutputFormat,
				Options: common.IgnoreParameterNotFoundOption,
			},
			resultFileName : "../test_files/all_resolved_no_paging_res.json",
		},
		{
			testCaseName : "Resolve all parameters with paging",
			numberOfParameters: 24,
			setting : common.Settings {
				InputFile: "../test_files/all_resolved_with_paging.json",
				OutputFile: testOutputFileName,
				OutputFormat: common.TextOutputFormat,
				Options: common.IgnoreParameterNotFoundOption,
			},
			resultFileName : "../test_files/all_resolved_with_paging_res.json",
		},
	}

	for i := 0; i < len(resolutionSettings); i++ {
		t.Logf("Running test: \"%s\"", resolutionSettings[i].testCaseName)

		serviceObject := NewServiceAllResolvedMockedObject(resolutionSettings[i].numberOfParameters)
		err := ResolveParametersInFile(&resolutionSettings[i].setting, &serviceObject)
		assert.Nil(t, err)

		expectedResultFileBuffer, _ := common.ReadTextFromFile(resolutionSettings[i].resultFileName)
		resultFileBuffer, _ := common.ReadTextFromFile(testOutputFileName)

		res := expectedResultFileBuffer == resultFileBuffer
		if res {
			t.Logf("\"%s\" succeeded %s", resolutionSettings[i].testCaseName, common.CheckMark)
		} else {
			t.Logf("\"%s\" has failed %s", resolutionSettings[i].testCaseName, common.BallotX)
		}

		assert.True(t, res)
	}
}
