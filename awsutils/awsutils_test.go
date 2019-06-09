// Package awsutils provides some helper function for common aws task.
package awsutils

import "testing"

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
