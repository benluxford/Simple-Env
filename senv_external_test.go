package senv_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/benluxford/Senv"
)

func TestLoad(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Errorf("unable to get the current working directory")
	}
	fileLocation := fmt.Sprintf("%s/.senv.test.file", workingDirectory)
	if err := senv.Load(fileLocation); err != nil {
		t.Fatalf("env failed to load %v", err)
	}
}

func TestGetSetVar(t *testing.T) {
	tt := []struct {
		key, value   string
		expectedPass bool
	}{
		{"a key", "some value", true},
		{"A_FORMATTED_KEY", "some other value", true},
		{"   key with spaces   ", "yet another value", true},
		{"you will not pass", "yet another value", false},
	}

	for _, tc := range tt {
		if tc.expectedPass {
			err := senv.SetVar(tc.key, tc.value)
			if err != nil {
				t.Errorf("unable to set key value pair, key:  %s, value: %s error: %v", tc.key, tc.value, err)
			}
		}
		value, err := senv.GetVar(tc.key)
		if err != nil && tc.expectedPass {
			t.Errorf("unable to get value for key: %s", tc.key)
		}
		if value != tc.value && tc.expectedPass {
			t.Errorf("cound not get correct value, expected: %s, got: %s", tc.value, value)
		}
	}
}
