// Package awsutils provides some helper function for common aws task.
package awsutils

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

func generateParamers(n int) map[string]string {
	parameters := make(map[string]string)

	for i := 1; i < n+1; i++ {
		index := strconv.Itoa(i)
		parameters["key"+index] = "value" + index
	}
	return parameters
}

func TestFindMissingParametresSuccess(t *testing.T) {

	value1 := "defValue1"
	value2 := "defValue2"
	requiredParam := map[string]*string{
		"key1": &value1,
		"key2": &value2,
		"key3": nil,
		"key4": nil,
	}
	parameters := generateParamers(4)
	err := findMissingParametres(requiredParam, parameters)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestFindMissingParametresFail(t *testing.T) {

	value1 := "defValue1"
	requiredParam := map[string]*string{
		"key1": &value1,
		"key2": nil,
		"key3": nil,
		"key4": nil,
	}
	parameters := map[string]string{
		"key2": "value2",
	}
	err := findMissingParametres(requiredParam, parameters)
	message := err.Error()

	if !strings.Contains(message, "key3") || !strings.Contains(message, "key4") {
		t.Errorf("Expected: key3 and key4, and got:  %s", message)
	}
}

func TestConvertToCfnParameter(t *testing.T) {

	parameters := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}
	cfnParam := convertToCfnParameter(parameters)

	if len(parameters) != len(cfnParam) {
		t.Errorf("Differnt number of parametres return ")
	}
}

func TestConvertToRequiredCfnParameter(t *testing.T) {

	requiredParam := map[string]*string{
		"key1": nil,
		"key2": nil,
	}
	parameters := generateParamers(4)
	cfnParam := convertToRequiredCfnParameter(requiredParam, parameters)
	if len(cfnParam) != 2 {
		t.Errorf("Two required parameters expected")
	}
}

func TestGetTeplateParameters(t *testing.T) {
	// Forgot to define client
	sError := Stack{
		Cfn: nil,
	}
	_, err := sError.GetTeplateParameters()

	if err.Error() != messageClientNotDefined {
		t.Errorf("Expected error :%s, and got %s", messageClientNotDefined, err.Error())
	}

	// Test success call
	mock := &mockedClient{
		RespValidateTemplateOutput: &cloudformation.ValidateTemplateOutput{
			Parameters: []*cloudformation.TemplateParameter{
				&cloudformation.TemplateParameter{ParameterKey: aws.String("key1")},
				&cloudformation.TemplateParameter{ParameterKey: aws.String("key2")}},
		},
	}
	s := Stack{
		Cfn: mock,
	}

	templateParam, err := s.GetTeplateParameters()
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(templateParam) != 2 {
		t.Errorf("Two parameters expected")
	}

}

func TestCreateStack(t *testing.T) {

	parameters := generateParamers(4)

	// Forgot to define client
	sError := Stack{
		Cfn: nil,
	}
	err := sError.CreateStack(parameters)

	if err.Error() != messageClientNotDefined {
		t.Errorf("Expected error :%s, and got %s", messageClientNotDefined, err.Error())
	}

	// Test success call
	mock := &mockedClient{
		RespValidateTemplateOutput: &cloudformation.ValidateTemplateOutput{
			Parameters: []*cloudformation.TemplateParameter{
				&cloudformation.TemplateParameter{ParameterKey: aws.String("key1")},
				&cloudformation.TemplateParameter{ParameterKey: aws.String("key2")}},
		},
	}
	s := Stack{
		Cfn: mock,
	}

	err = s.CreateStack(parameters)
	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestCreateChangeSet(t *testing.T) {
	parameters := generateParamers(4)
	// Forgot to define client
	sError := Stack{
		Cfn: nil,
	}
	err := sError.CreateChangeSet(parameters)

	if err.Error() != messageClientNotDefined {
		t.Errorf("Expected error :%s, and got %s", messageClientNotDefined, err.Error())
	}

	// Test success call
	mock := &mockedClient{
		RespValidateTemplateOutput: &cloudformation.ValidateTemplateOutput{
			Parameters: []*cloudformation.TemplateParameter{
				&cloudformation.TemplateParameter{ParameterKey: aws.String("key1")},
				&cloudformation.TemplateParameter{ParameterKey: aws.String("key2")}},
		},
	}
	s := Stack{
		Cfn: mock,
	}

	err = s.CreateChangeSet(parameters)
	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestCreateOrUpdate(t *testing.T) {
	parameters := generateParamers(4)
	// Forgot to define client
	sError := Stack{
		Cfn: nil,
	}
	err := sError.CreateOrUpdate(parameters)

	if err.Error() != messageClientNotDefined {
		t.Errorf("Expected error :%s, and got %s", messageClientNotDefined, err.Error())
	}

	// Test success call
	mock := &mockedClient{
		RespValidateTemplateOutput: &cloudformation.ValidateTemplateOutput{
			Parameters: []*cloudformation.TemplateParameter{
				&cloudformation.TemplateParameter{ParameterKey: aws.String("key1")},
				&cloudformation.TemplateParameter{ParameterKey: aws.String("key2")}},
		},
	}
	s := Stack{
		Cfn: mock,
	}

	err = s.CreateOrUpdate(parameters)
	if err != nil {
		t.Errorf(err.Error())
	}

}

/*Mock stuff*/

type mockedClient struct {
	cloudformationiface.CloudFormationAPI
	RespValidateTemplateOutput *cloudformation.ValidateTemplateOutput
}

func (m *mockedClient) ValidateTemplate(in *cloudformation.ValidateTemplateInput) (*cloudformation.ValidateTemplateOutput, error) {
	return m.RespValidateTemplateOutput, nil
}

func (m *mockedClient) DescribeStacks(in *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	return nil, fmt.Errorf("Not found error")
}

func (m *mockedClient) CreateStack(in *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return &cloudformation.CreateStackOutput{}, nil
}

func (m *mockedClient) CreateChangeSet(in *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
	return &cloudformation.CreateChangeSetOutput{}, nil
}

func (m *mockedClient) WaitUntilStackCreateComplete(in *cloudformation.DescribeStacksInput) error {
	return nil
}
