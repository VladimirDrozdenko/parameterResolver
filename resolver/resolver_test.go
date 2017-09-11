package resolver

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
)

func TestExtractParametersFromText(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}."
	resolvedParameters, err := ExtractParametersFromText(&serviceObject, text, ResolveOptions{
		IgnoreSecureParameters: false,
	})

	expectedResult := map[string]SsmParameterInfo {
		"ssm:/a/b/c/param1": {Name: "/a/b/c/param1", Type: stringType, Value: "value_/a/b/c/param1"},
		"ssm-secure:param2": {Name: "param2", Type: secureStringType, Value: "value_param2"},
	}

	assert.Nil(t, err)
	assert.NotNil(t, resolvedParameters)
	assert.True(t, reflect.DeepEqual(resolvedParameters, expectedResult))
}

func TestExtractParametersFromTextIgnoreSecureParams(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}."
	resolvedParameters, err := ExtractParametersFromText(&serviceObject, text, ResolveOptions{
		IgnoreSecureParameters: true,
	})

	expectedResult := map[string]SsmParameterInfo {
		"ssm:/a/b/c/param1": {Name: "/a/b/c/param1", Type: stringType, Value: "value_/a/b/c/param1"},
	}

	assert.Nil(t, err)
	assert.NotNil(t, resolvedParameters)
	assert.True(t, reflect.DeepEqual(resolvedParameters, expectedResult))
}

func TestResolveParameterReferenceList(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	parameterReferences := []string {
		"ssm:param1",
		"ssm:param2",
		"ssm-secure:/a/b/param1",
		"ssm-secure:param4",
	}

	resolvedParameters, err := ResolveParameterReferenceList(&serviceObject, parameterReferences, ResolveOptions{
		IgnoreSecureParameters: false,
	})

	expectedResult := map[string]SsmParameterInfo {
		"ssm:param1": {Name: "param1", Type: stringType, Value: "value_param1"},
		"ssm:param2": {Name: "param2", Type: stringType, Value: "value_param2"},
		"ssm-secure:/a/b/param1": {Name: "/a/b/param1", Type: secureStringType, Value: "value_/a/b/param1"},
		"ssm-secure:param4": {Name: "param4", Type: secureStringType, Value: "value_param4"},
	}

	assert.Nil(t, err)
	assert.NotNil(t, resolvedParameters)
	assert.True(t, reflect.DeepEqual(resolvedParameters, expectedResult))
}

func TestResolveParameterReferenceListIgnoreSecureParams(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	parameterReferences := []string {
		"ssm:param1",
		"ssm:param2",
		"ssm-secure:/a/b/param1",
		"ssm-secure:param4",
	}

	resolvedParameters, err := ResolveParameterReferenceList(&serviceObject, parameterReferences, ResolveOptions{
		IgnoreSecureParameters: true,
	})

	expectedResult := map[string]SsmParameterInfo {
		"ssm:param1": {Name: "param1", Type: stringType, Value: "value_param1"},
		"ssm:param2": {Name: "param2", Type: stringType, Value: "value_param2"},
	}

	assert.Nil(t, err)
	assert.NotNil(t, resolvedParameters)
	assert.True(t, reflect.DeepEqual(resolvedParameters, expectedResult))
}

func TestParseParametersFromTextIntoMapSecureAllowed(t *testing.T) {
	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}, {{ ssm:/a/b/c/param1  }}."
	expectedList := []string {"ssm:/a/b/c/param1", "ssm-secure:param2"}

	list, err := parseParametersFromTextIntoMap(text)

	assert.Nil(t, err)
	assert.NotNil(t, list)

	sort.Slice(expectedList, func(i, j int) bool { return expectedList[i] < expectedList[j] })
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	assert.True(t, reflect.DeepEqual(list, expectedList))
}

func TestResolveParametersInText(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}."
	output, err := ResolveParametersInText(&serviceObject, text, ResolveOptions{
		IgnoreSecureParameters: false,
	})

	expectedOutput := `Some text value_/a/b/c/param1, some more text value_param2.`

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.True(t, expectedOutput == output)
}

func TestResolveParametersInTextIgnoreSecureParams(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}."
	output, err := ResolveParametersInText(&serviceObject, text, ResolveOptions{
		IgnoreSecureParameters: true,
	})

	expectedOutput := `Some text value_/a/b/c/param1, some more text {{ssm-secure:param2}}.`

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.True(t, expectedOutput == output)
}
