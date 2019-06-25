package awsutils

import (
	"testing"
)

func TestGeneratePasswordErrors(t *testing.T) {

	empty := ""
	_, err := GeneratePassword(0, 1, 1, nil, nil, nil)
	if err.Error() != exceedsTotalLength {
		t.Errorf("expected %q to be %q", err, exceedsTotalLength)
	}

	_, err = GeneratePassword(0, 1, 1, nil, &empty, nil)
	if err.Error() != allowedNumbersMustBeDefined {
		t.Errorf("expected %q to be %q", err, allowedNumbersMustBeDefined)
	}

	_, err = GeneratePassword(0, 1, 1, nil, nil, &empty)
	if err.Error() != allowedSymbolsMustBeDefined {
		t.Errorf("expected %q to be %q", err, allowedNumbersMustBeDefined)
	}
}

func TestGeneratePassword(t *testing.T) {

	letters := "abAB"
	numbers := "01"
	symbols := "!?"
	pass, err := GeneratePassword(4, 1, 1, &letters, &numbers, &symbols)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(pass) != 4 {
		t.Errorf("the password was expected to be of size four")
	}

}
