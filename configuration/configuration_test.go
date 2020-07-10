package configuration

import "testing"

func Test_profile_name_default_is_invalid(t *testing.T) {

	isValid, _ := validateProfile("default")

	if !isValid {
		t.Errorf("profile name 'default' must be valid")
	}
}

func Test_profile_name_with_specialcharacter_is_invalid(t *testing.T) {

	isValid, _ := validateProfile("default!@#$%^&*()+=?></.,';\":`~")

	if isValid {
		t.Errorf("profile name with special characters should be invalid")
	}
}

func Test_profile_name_with_allowed_characters_is_invali(t *testing.T) {

	isValid, _ := validateProfile("1234567890qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM-_.")

	if !isValid {
		t.Errorf("profile name with latin chars numbers and -_. should be valid")
	}
}
