package awsutils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

const (
	messageClientNotDefined = "Aws Client not defined"
)

//Stack ... Aws Cloud formation stack
type Stack struct {
	Cfn          cloudformationiface.CloudFormationAPI
	Name         string
	TemplateURL  string
	Capabilities []string
	Status       *string
}

func (s *Stack) InitilizeCfn(region string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	s.Cfn = cloudformation.New(sess)
}

//CreateOrUpdate ... creates a stack or creates a change set for an existing stack based on given parameters
func (s *Stack) CreateOrUpdate(parameters map[string]string) error {

	if s.Cfn == nil {
		return fmt.Errorf(messageClientNotDefined)
	}

	templateParam, err := s.getTeplateParameters()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := findMissingParametres(templateParam, parameters); err != nil {
		log.Println(err.Error())
		return err
	}

	cfnParameters := convertToRequiredCfnParameter(templateParam, parameters)
	input := cloudformation.DescribeStacksInput{StackName: &s.Name}
	_, err = s.Cfn.DescribeStacks(&input)

	if err != nil {
		err = s.createStack(cfnParameters)
	} else {
		err = s.createChangeSet(cfnParameters)
	}
	return err
}
func findMissingParametres(templateParam map[string]*string, parameters map[string]string) error {
	missing := make([]string, 0)
	for key, defaultValue := range templateParam {
		_, doesKeyExist := parameters[key]
		if !doesKeyExist && defaultValue == nil {
			missing = append(missing, key)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("Missing: [%s]", strings.Join(missing, ","))
}
func convertToCfnParameter(parameters map[string]string) []*cloudformation.Parameter {
	result := make([]*cloudformation.Parameter, 0)
	for key, value := range parameters {
		result = append(result, &cloudformation.Parameter{
			ParameterKey:   aws.String(key),
			ParameterValue: aws.String(value),
		})
	}
	return result
}
func convertToRequiredCfnParameter(templateParam map[string]*string, parameters map[string]string) []*cloudformation.Parameter {
	result := make([]*cloudformation.Parameter, 0)
	for key := range templateParam {
		value, ok := parameters[key]
		if ok {
			result = append(result, &cloudformation.Parameter{
				ParameterKey:   aws.String(key),
				ParameterValue: aws.String(value),
			})
		}
	}
	return result
}

//ReadOutputs ...
func (s *Stack) ReadOutputs() (map[string]string, error) {
	if s.Cfn == nil {
		return nil, fmt.Errorf(messageClientNotDefined)
	}
	parameters := make(map[string]string)
	input := cloudformation.DescribeStacksInput{StackName: &s.Name}

	res, err := s.Cfn.DescribeStacks(&input)
	if err != nil {
		return nil, err
	}
	for _, stack := range res.Stacks {
		for _, output := range stack.Outputs {
			parameters[*output.OutputKey] = *output.OutputValue
		}
	}
	return parameters, nil
}

//LoadParameters ...
func LoadParameters(fileName string) (map[string]string, error) {
	parameters := make(map[string]string)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words := strings.Split(scanner.Text(), "=")
		key := words[0]
		value := words[1]
		parameters[key] = value
	}
	return parameters, scanner.Err()
}

//LoadEnvironmentVariables ...
func LoadEnvironmentVariables() (map[string]string, error) {

	parameters := make(map[string]string)
	for _, pair := range os.Environ() {

		keyValues := strings.Split(pair, "=")
		key := keyValues[0]
		value := keyValues[1]
		parameters[key] = value
	}
	return parameters, nil
}

//GetAllStacksBy ...
func GetAllStacksBy(region string) ([]Stack, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	svc := cloudformation.New(sess)

	var filter = []*string{
		aws.String("CREATE_IN_PROGRESS"),
		aws.String("CREATE_FAILED"),
		aws.String("CREATE_COMPLETE"),
		aws.String("ROLLBACK_IN_PROGRESS"),
		aws.String("ROLLBACK_FAILED"),
		aws.String("ROLLBACK_COMPLETE"),
		aws.String("DELETE_IN_PROGRESS"),
		aws.String("DELETE_FAILED"),
		//aws.String("DELETE_COMPLETE"),
		aws.String("UPDATE_IN_PROGRESS"),
		aws.String("UPDATE_COMPLETE_CLEANUP_IN_PROGRESS"),
		aws.String("UPDATE_COMPLETE"),
		aws.String("UPDATE_ROLLBACK_IN_PROGRESS"),
		aws.String("UPDATE_ROLLBACK_FAILED"),
		aws.String("UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS"),
		aws.String("UPDATE_ROLLBACK_COMPLETE"),
		aws.String("REVIEW_IN_PROGRESS")}
	input := &cloudformation.ListStacksInput{StackStatusFilter: filter}

	resp, err := svc.ListStacks(input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	results := make([]Stack, 0)

	for _, summary := range resp.StackSummaries {
		results = append(results, Stack{Name: *summary.StackName, Status: summary.StackStatus})
	}
	return results, nil
}

//GetTeplateParameters ...
func (s *Stack) GetTeplateParameters() (map[string]*string, error) {
	if s.Cfn == nil {
		return nil, fmt.Errorf(messageClientNotDefined)
	}
	return s.getTeplateParameters()
}
func (s *Stack) getTeplateParameters() (map[string]*string, error) {

	input := &cloudformation.ValidateTemplateInput{TemplateURL: &s.TemplateURL}
	resp, err := s.Cfn.ValidateTemplate(input)
	if err != nil {
		return nil, err
	}
	resultParameters := make(map[string]*string)
	for _, tp := range resp.Parameters {
		resultParameters[*tp.ParameterKey] = tp.DefaultValue
	}
	return resultParameters, nil
}

//CreateStack ...
func (s *Stack) CreateStack(parameters map[string]string) error {
	if s.Cfn == nil {
		return fmt.Errorf(messageClientNotDefined)
	}
	cfnParameters := convertToCfnParameter(parameters)
	return s.createStack(cfnParameters)
}
func (s *Stack) createStack(parameters []*cloudformation.Parameter) error {
	input := &cloudformation.CreateStackInput{
		TemplateURL:  aws.String(s.TemplateURL),
		StackName:    aws.String(s.Name),
		Capabilities: aws.StringSlice(s.Capabilities),
		Parameters:   parameters}

	_, err := s.Cfn.CreateStack(input)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Wait until stack is created
	desInput := &cloudformation.DescribeStacksInput{StackName: aws.String(s.Name)}
	err = s.Cfn.WaitUntilStackCreateComplete(desInput)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//CreateChangeSet ...
func (s *Stack) CreateChangeSet(parameters map[string]string) error {
	if s.Cfn == nil {
		return fmt.Errorf(messageClientNotDefined)
	}
	cfnParameters := convertToCfnParameter(parameters)
	return s.createChangeSet(cfnParameters)
}
func (s *Stack) createChangeSet(parameters []*cloudformation.Parameter) error {

	t := time.Now()
	changeSetName := s.Name + "-" + t.Format("20060102030405")
	input := &cloudformation.CreateChangeSetInput{
		TemplateURL:   aws.String(s.TemplateURL),
		StackName:     aws.String(s.Name),
		ChangeSetName: aws.String(changeSetName),
		Parameters:    parameters}

	_, err := s.Cfn.CreateChangeSet(input)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Wait until stack is created
	desInput := &cloudformation.DescribeStacksInput{StackName: aws.String(s.Name)}
	err = s.Cfn.WaitUntilStackCreateComplete(desInput)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
