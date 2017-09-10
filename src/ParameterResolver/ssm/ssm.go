package ssm

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
	"parameterResolver/common"
	"strconv"
)

type ISsmParameterService interface {
	GetParameters(names []string) (map[string]string, []string, error)
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
		log.Println("no explict region configuration. So now retriving ec2metadata...")
		region, err := ec2metadata.New(currentSession).Region()
		if err != nil {
			return nil, err
		}
		currentSession.Config.Region = aws.String(region)
	}

	log.Printf("Region: %s\n", *currentSession.Config.Region)

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

func (s *Service) GetParameters(names []string) (map[string]string, []string, error) {

	log.Println("Making a call to SSM Parameters Store to fetch " + strconv.Itoa(len(names)) + " parameter(s)")
	parametersOutput, err := s.SSMClient.GetParameters(&ssm.GetParametersInput{
		Names:          aws.StringSlice(names),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, nil, err
	}

	parameters := make(map[string]string)

	for i := 0; i < len(parametersOutput.Parameters); i++ {
		param := parametersOutput.Parameters[i]
		parameters[*param.Name] = *param.Value
	}

	var invalidParameters []string
	for i := 0; i < len(parametersOutput.InvalidParameters); i++ {
		invalidParameters = append(invalidParameters, *parametersOutput.InvalidParameters[i])
	}

	return parameters, invalidParameters, nil
}

func GetParametersFromSsmParameterStore(s ISsmParameterService, namesToFetch map[string]string, options string) (map[string]string, error) {

	parameters := make(map[string]string)
	for k,v := range namesToFetch {
		parameters[k] = v
	}

	var totalParams = len(namesToFetch)
	for totalParams > 0 {

		var paramsBatch []string

		count := 0
		for paramPair := range namesToFetch {
			paramsBatch = append(paramsBatch, paramPair)
			delete(namesToFetch, paramPair)

			totalParams--
			count++
			if count == common.MaxParametersRetrievedFromSsm {
				break
			}
		}

		results, invalidParameters, err := s.GetParameters(paramsBatch)
		if err != nil {
			return nil, err
		}

		if options == common.FailOnParameterNotFoundOption && invalidParameters != nil {
			return nil, errors.New("The following parameter(s) not found: " + strings.Join(invalidParameters, ","))
		}

		for name, value := range results {
			parameters[name] = value
		}
	}

	return parameters, nil
}

