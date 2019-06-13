// Package awsutils provides some helper function for common aws task.
package awsutils

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

func TestFindMissingParametresSuccess(t *testing.T) {

	value1 := "defValue1"
	value2 := "defValue2"
	requiredParam := map[string]*string{
		"key1": &value1,
		"key2": &value2,
		"key3": nil,
		"key4": nil,
	}
	parameters := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}
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
	expectedError := "Missing: [key3,key4]"
	if message != expectedError {
		t.Errorf("Expected error: %s", expectedError)
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
	parameters := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}
	cfnParam := convertToRequiredCfnParameter(requiredParam, parameters)
	if len(cfnParam) != 2 {
		t.Errorf("Two required parameters expected")
	}
}

type mockedClient struct {
	cloudformationiface.CloudFormationAPI
	RespValidateTemplateOutput *cloudformation.ValidateTemplateOutput
}

func (m *mockedClient) ValidateTemplate(in *cloudformation.ValidateTemplateInput) (*cloudformation.ValidateTemplateOutput, error) {
	// Only need to return mocked response output
	return m.RespValidateTemplateOutput, nil
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
				&cloudformation.TemplateParameter{ParameterKey: aws.String("Key1")},
				&cloudformation.TemplateParameter{ParameterKey: aws.String("Key2")}},
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
