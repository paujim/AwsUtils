package awsutils

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type Store struct {
	ssmClient ssmiface.SSMAPI
}

func NewStore(client ssmiface.SSMAPI) Store {
	return Store{ssmClient: client}
}

func (s *Store) GetParameter(keyname string) (*string, error) {
	if s.ssmClient == nil {
		return nil, errors.New(messageClientNotDefined)
	}

	withDecryption := true
	input := &ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	}
	param, err := s.ssmClient.GetParameter(input)
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
	_, err := s.ssmClient.PutParameter(input)
	return err
}
