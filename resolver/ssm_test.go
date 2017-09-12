package resolver

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ServiceMockedObjectWithRecords struct {
	ISsmParameterService
	records map[string]SsmParameterInfo
}

func NewServiceMockedObjectWithExtraRecords(
	records map[string]SsmParameterInfo) ServiceMockedObjectWithRecords {
	return ServiceMockedObjectWithRecords{
		records: records,
	}
}

func (m *ServiceMockedObjectWithRecords) callGetParameters(parameterReferences []string) (map[string]SsmParameterInfo, error) {
	parameters := make(map[string]SsmParameterInfo)

	for i := 0; i < len(parameterReferences); i++ {

		value, contains := m.records[parameterReferences[i]]
		if !contains {
			return nil, errors.New("error: " + parameterReferences[i] + " cannot be resolved")
		}

		parameters[parameterReferences[i]] = value
	}

	return parameters, nil
}

func TestGetParametersFromSsmParameterStoreWithAllResolvedNoPaging(t *testing.T) {
	parametersList := []string{}
	expectedValues := map[string]SsmParameterInfo{}

	for i := 0; i < maxParametersRetrievedFromSsm/2; i++ {
		name := "name_" + strconv.Itoa(i)
		key := ssmNonSecurePrefix + name
		parametersList = append(parametersList, key)

		expectedValues[key] = SsmParameterInfo{
			Name:  name,
			Value: "value_" + name,
			Type:  "String",
		}
	}

	serviceObject := NewServiceMockedObjectWithExtraRecords(expectedValues)

	t.Log("Testing getParametersFromSsmParameterStore API for all parameters present without paging...")
	retrievedValues, err := getParametersFromSsmParameterStore(&serviceObject, parametersList)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expectedValues, retrievedValues))
}

func TestGetParametersFromSsmParameterStoreWithAllResolvedWithPaging(t *testing.T) {
	parametersList := []string{}
	expectedValues := map[string]SsmParameterInfo{}

	for i := 0; i < maxParametersRetrievedFromSsm/5; i++ {
		name := "name_" + strconv.Itoa(i)
		key := ssmSecurePrefix + name
		parametersList = append(parametersList, key)

		expectedValues[key] = SsmParameterInfo{
			Name:  name,
			Value: "value_" + name,
			Type:  secureStringType,
		}
	}

	serviceObject := NewServiceMockedObjectWithExtraRecords(expectedValues)

	t.Log("Testing getParametersFromSsmParameterStore API for all parameters present with paging...")
	retrievedValues, err := getParametersFromSsmParameterStore(&serviceObject, parametersList)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expectedValues, retrievedValues))
}

func TestGetParametersFromSsmParameterStoreWithUnresolvedIgnoreNoPaging(t *testing.T) {
	parametersList := []string{}
	for i := 0; i < 2; i++ {
		key := "{{ssm:name_" + strconv.Itoa(i) + "}}"
		parametersList = append(parametersList, key)
	}

	serviceObject := NewServiceMockedObjectWithExtraRecords(map[string]SsmParameterInfo{})

	t.Log("Testing getParametersFromSsmParameterStore API for all unresolved parameters...")
	_, err := getParametersFromSsmParameterStore(&serviceObject, parametersList)
	assert.NotNil(t, err)
}
