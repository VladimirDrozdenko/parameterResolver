package ssm

import (
	"testing"
	"parameterResolver/common"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
)

const unresolvedKeyName = "name_0"

type ServiceAllResolvedMockedObject struct {
	ISsmParameterService
}

type ServiceWithUnresolvedMockedObject struct {
	ISsmParameterService
}

func (m *ServiceAllResolvedMockedObject) GetParameters(names []string) (map[string]string, []string, error) {

	parameters := make(map[string]string)

	for i := 0; i < len(names); i++ {
		parameters[names[i]] = "value_" + names[i]
	}

	return parameters, nil, nil
}

func (m *ServiceWithUnresolvedMockedObject) GetParameters(names []string) (map[string]string, []string, error) {

	parameters := make(map[string]string)

	var invalidParameters []string

	for i := 0; i < len(names); i++ {
		if names[i] == unresolvedKeyName {
			invalidParameters = append(invalidParameters, names[i])
		} else {
			parameters[names[i]] = "value_" + names[i]
		}
	}

	return parameters, invalidParameters, nil
}


func TestGetParametersFromSsmParameterStoreWithAllResolvedNoPaging(t *testing.T) {
	serviceObject := new(ServiceAllResolvedMockedObject)

	parametersMap := map[string]string {}
	for i := 0; i < common.MaxParametersRetrievedFromSsm / 2; i++ {
		name := "name_" + strconv.Itoa(i)
		parametersMap[name] = "{{ssm:" + name + "}}"
	}

	expectedValues := map[string]string{}
	for key, _ := range parametersMap {
		expectedValues[key] = "value_" + key
	}

	t.Log("Testing GetParametersFromSsmParameterStore API for all parameters present without paging...")
	retrievedValues, err := GetParametersFromSsmParameterStore(serviceObject, parametersMap, common.FailOnParameterNotFoundOption)
	assert.Nil(t, err, "aaa")
	assert.True(t, reflect.DeepEqual(expectedValues, retrievedValues))
}


func TestGetParametersFromSsmParameterStoreWithAllResolvedWithPaging(t *testing.T) {
	serviceObject := new(ServiceAllResolvedMockedObject)

	parametersMap := map[string]string {}
	for i := 0; i < common.MaxParametersRetrievedFromSsm + 5; i++ {
		name := "name_" + strconv.Itoa(i)
		parametersMap[name] = "{{ssm:" + name + "}}"
	}

	expectedValues := map[string]string{}
	for key, _ := range parametersMap {
		expectedValues[key] = "value_" + key
	}

	retrievedValues, err := GetParametersFromSsmParameterStore(serviceObject, parametersMap, common.FailOnParameterNotFoundOption)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expectedValues, retrievedValues))
}


func TestGetParametersFromSsmParameterStoreWithUnresolvedIgnoreNoPaging(t *testing.T) {
	serviceObject := new(ServiceWithUnresolvedMockedObject)

	parametersMap := map[string]string {}
	for i := 0; i < common.MaxParametersRetrievedFromSsm / 2; i++ {
		name := "name_" + strconv.Itoa(i)
		parametersMap[name] = "{{ssm:" + name + "}}"
	}

	expectedValues := map[string]string{}

	for key, _ := range parametersMap {
		if key == unresolvedKeyName {
			expectedValues[key] = "{{ssm:" + key + "}}"
		} else {
			expectedValues[key] = "value_" + key
		}
	}

	retrievedValues, err := GetParametersFromSsmParameterStore(serviceObject, parametersMap, common.IgnoreParameterNotFoundOption)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expectedValues, retrievedValues))
}


func TestGetParametersFromSsmParameterStoreWithUnresolvedIgnoreWithPaging(t *testing.T) {
	serviceObject := new(ServiceWithUnresolvedMockedObject)

	parametersMap := map[string]string {}
	for i := 0; i < common.MaxParametersRetrievedFromSsm + 5; i++ {
		name := "name_" + strconv.Itoa(i)
		parametersMap[name] = "{{ssm:" + name + "}}"
	}

	expectedValues := map[string]string{}

	for key, _ := range parametersMap {
		if key == unresolvedKeyName {
			expectedValues[key] = "{{ssm:" + key + "}}"
		} else {
			expectedValues[key] = "value_" + key
		}
	}

	retrievedValues, err := GetParametersFromSsmParameterStore(serviceObject, parametersMap, common.IgnoreParameterNotFoundOption)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expectedValues, retrievedValues))
}
