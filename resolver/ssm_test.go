package resolver

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"errors"
	"strings"
)

type ServiceMockedObject struct {
	ISsmParameterService
	generateUnresolved bool
}

func NewServiceMockedObject(genUnresolved bool) ServiceMockedObject {
	return ServiceMockedObject {
		generateUnresolved: genUnresolved,
	}
}

func (m *ServiceMockedObject) callGetParameters(parameterReferences []string) (map[string]SsmParameterInfo, error) {
	parameters := make(map[string]SsmParameterInfo)

	for i := 0; i < len(parameterReferences); i++ {
		key := extractParameterNameFromReference(parameterReferences[i])

		paramType := stringType
		if strings.HasPrefix(parameterReferences[i], ssmSecurePrefix) {
			paramType = secureStringType
		}
		parameters[parameterReferences[i]] = SsmParameterInfo {
			Name: key,
			Value: "value_" + key,
			Type: paramType,
		}
	}

	if m.generateUnresolved {
		return nil, errors.New("error")
	}

	return parameters, nil
}

func TestGetParametersFromSsmParameterStoreWithAllResolvedNoPaging(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	parametersList := []string {}
	expectedValues := map[string]SsmParameterInfo {}

	for i := 0; i < maxParametersRetrievedFromSsm / 2; i++ {
		name := "name_" + strconv.Itoa(i)
		key := ssmNonSecurePrefix + name
		parametersList = append(parametersList, key)

		expectedValues[key] = SsmParameterInfo {
			Name: name,
			Value: "value_" + name,
			Type: "String",
		}
	}

	t.Log("Testing getParametersFromSsmParameterStore API for all parameters present without paging...")
	retrievedValues, err := getParametersFromSsmParameterStore(&serviceObject, parametersList)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expectedValues, retrievedValues))
}


func TestGetParametersFromSsmParameterStoreWithAllResolvedWithPaging(t *testing.T) {
	serviceObject := NewServiceMockedObject(false)

	parametersList := []string {}
	expectedValues := map[string]SsmParameterInfo {}

	for i := 0; i < maxParametersRetrievedFromSsm / 5; i++ {
		name := "name_" + strconv.Itoa(i)
		key := ssmSecurePrefix + name
		parametersList = append(parametersList, key)

		expectedValues[key] = SsmParameterInfo {
			Name: name,
			Value: "value_" + name,
			Type: secureStringType,
		}
	}

	t.Log("Testing getParametersFromSsmParameterStore API for all parameters present with paging...")
	retrievedValues, err := getParametersFromSsmParameterStore(&serviceObject, parametersList)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expectedValues, retrievedValues))
}


func TestGetParametersFromSsmParameterStoreWithUnresolvedIgnoreNoPaging(t *testing.T) {
	serviceObject := NewServiceMockedObject(true)

	parametersList := []string {}
	for i := 0; i < 2; i++ {
		key := "{{ssm:name_" + strconv.Itoa(i) + "}}"
		parametersList = append(parametersList, key)
	}

	t.Log("Testing getParametersFromSsmParameterStore API for all unresolved parameters...")
	_, err := getParametersFromSsmParameterStore(&serviceObject, parametersList)
	assert.NotNil(t, err)
}
