package resolver

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"errors"
	"strings"
//	"strconv"
)

//
// Maximum number of parameters that can be requested from SSM Parameter store in one GetParameters request
const maxParametersRetrievedFromSsm = 10


type ISsmParameterService interface {
	callGetParameters(parameterReferences []string) (map[string]SsmParameterInfo, error)
}

type Service struct {
	SSMClient *ssm.SSM
}

func NewService() (service *Service, err error) {
	currentSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return
	}

	if *currentSession.Config.Region == "" {
		log.Println("There is no explict region configuration, retriving ec2metadata...")
		region, err := ec2metadata.New(currentSession).Region()
		if err != nil {
			return nil, err
		}
		currentSession.Config.Region = aws.String(region)
	}

	//log.Printf("Region: %s\n", *currentSession.Config.Region)

	var client *ssm.SSM
	if arn := os.Getenv("SSM2ENV_ASSUME_ROLE_ARN"); arn != "" {
		credentials := stscreds.NewCredentials(currentSession, arn)
		client = ssm.New(currentSession, &aws.Config{Credentials: credentials})
	} else {
		client = ssm.New(currentSession)
	}

	service = &Service {
		SSMClient: client,
	}

	return
}

//
// This function takes a list of at most maxParametersRetrievedFromSsm(=10) ssm parameter name references like (ssm:name).
// It returns a map<param-ref, SsmParameterInfo>.
func (s *Service) callGetParameters(parameterReferences []string) (map[string]SsmParameterInfo, error) {

	//log.Println("Making a call to SSM Parameters Store to fetch " + strconv.Itoa(len(parameterReferences)) + " parameter(s)")

	ref2NameMapper := make(map[string]string)

	for i := 0; i < len(parameterReferences); i++ {
		nameWithoutPrefix := extractParameterNameFromReference(parameterReferences[i])
		ref2NameMapper[nameWithoutPrefix] = parameterReferences[i]
		parameterReferences[i] = nameWithoutPrefix
	}

	parametersOutput, err := s.SSMClient.GetParameters(&ssm.GetParametersInput{
		Names:          aws.StringSlice(parameterReferences),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	if len(parametersOutput.InvalidParameters) > 0 {
		invalidParameters := []string{}
		for _, p := range parametersOutput.InvalidParameters {
			invalidParameters = append(invalidParameters, *p)
		}
		return nil, errors.New("The following parameter(s) cannot be resolved: " + strings.Join(invalidParameters, ","))
	}

	resolvedParametersMap := map[string]SsmParameterInfo{}
	for i := 0; i < len(parametersOutput.Parameters); i++ {
		param := parametersOutput.Parameters[i]
		resolvedParametersMap[ref2NameMapper[*param.Name]] = SsmParameterInfo {
			Name: *param.Name,
			Type: *param.Type,
			Value: *param.Value,
		}
	}

	return resolvedParametersMap, nil
}

//
// This function takes as an input a list of references to the SSMParameterService and return a map <reference, SSMParameterInfo>
func getParametersFromSsmParameterStore(s ISsmParameterService, parametersToFetch []string) (map[string]SsmParameterInfo, error) {

	outputMap := make(map[string]SsmParameterInfo)

	var totalParams = len(parametersToFetch)
	var startPos = 0
	for totalParams > 0 {

		var paramsBatch []string

		var count = 0
		for i := startPos; i < len(parametersToFetch) && count < maxParametersRetrievedFromSsm; i++ {
			paramsBatch = append(paramsBatch, parametersToFetch[i])

			totalParams--
			count++
			startPos++
		}

		results, err := s.callGetParameters(paramsBatch)
		if err != nil {
			return nil, err
		}

		for name, value := range results {
			outputMap[name] = value
		}
	}

	return outputMap, nil
}

func extractParameterNameFromReference(parameterReference string) string {
	return parameterReference[strings.Index(parameterReference, ":") + 1 :]
}