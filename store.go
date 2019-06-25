package awsutils

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type Store struct {
	client ssmiface.SSMAPI
}

func (s *Store) GetParameter(keyname string) (*string, error) {
	if s == nil {
		return nil, errors.New(messageClientNotDefined)
	}

	withDecryption := true
	input := &ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	}
	param, err := s.client.GetParameter(input)
	if err != nil {
		return nil, err
	}
	return param.Parameter.Value, nil
}

func (s *Store) PutParameter(keyname, value string) error {
	if s == nil {
		return nil
	}
	input := &ssm.PutParameterInput{
		Name:  aws.String(keyname),
		Value: aws.String(value),
		Type:  aws.String(ssm.ParameterTypeSecureString),
	}
	_, err := s.client.PutParameter(input)
	return err
}
