package resolver

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

func TestExtractParametersFromText(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}."
	resolvedParameters, err := ExtractParametersFromText(&serviceObject, text, ResolveOptions{
		ResolveSecureParameters: true,
	})

	expectedResult := map[string]SsmParameterInfo {
		"ssm:/a/b/c/param1": {Name: "/a/b/c/param1", Type: stringType, Value: "value_/a/b/c/param1"},
		"ssm-secure:param2": {Name: "param2", Type: secureStringType, Value: "value_param2"},
	}

	assert.Nil(t, err)
	assert.NotNil(t, resolvedParameters)
	assert.True(t, reflect.DeepEqual(resolvedParameters, expectedResult))
}

func TestExtractParametersFromTextNoSecureParams(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}."
	resolvedParameters, err := ExtractParametersFromText(&serviceObject, text, ResolveOptions{
		ResolveSecureParameters: false,
	})

	assert.NotNil(t, err)
	assert.Nil(t, resolvedParameters)
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
		ResolveSecureParameters: true,
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

func TestResolveParameterReferenceListNoSecureParams(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	parameterReferences := []string {
		"ssm:param1",
		"ssm:param2",
		"ssm-secure:/a/b/param1",
		"ssm-secure:param4",
	}

	resolvedParameters, err := ResolveParameterReferenceList(&serviceObject, parameterReferences, ResolveOptions{
		ResolveSecureParameters: false,
	})

	assert.NotNil(t, err)
	assert.Nil(t, resolvedParameters)
}

func TestParseParametersFromTextIntoMapSecureAllowed(t *testing.T) {
	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}."
	expectedList := []string {"ssm:/a/b/c/param1", "ssm-secure:param2"}

	list, err := parseParametersFromTextIntoMap(text, ResolveOptions{
		ResolveSecureParameters: true,
	})

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.True(t, reflect.DeepEqual(list, expectedList))
}

func TestParseParametersFromTextIntoMapSecureNotAllowed(t *testing.T) {
	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}...."

	list, err := parseParametersFromTextIntoMap(text, ResolveOptions{
		ResolveSecureParameters: false,
	})

	assert.NotNil(t, err)
	assert.Nil(t, list)
}

func TestResolveParametersInText(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	text := "Some text {{ ssm:/a/b/c/param1}}, some more text {{ssm-secure:param2}}."
	output, err := ResolveParametersInText(&serviceObject, text, ResolveOptions{
		ResolveSecureParameters: true,
	})

	expectedOutput := `Some text value_/a/b/c/param1, some more text value_param2.`

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.True(t, expectedOutput == output)
}
