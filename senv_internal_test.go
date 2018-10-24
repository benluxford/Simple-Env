package senv

import (
	"fmt"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Errorf("unable to get the current working directory")
	}
	tt := []struct {
		fileName   string
		key, value string
		expected   bool
	}{
		{fmt.Sprintf("%s/.senv.test.file", workingDirectory), "something", "else", true},
		{fmt.Sprintf("%s/.senv.test.file", workingDirectory), "package", "Senv", true},
		{"file that doesn't exist", "package", "Senv", false},
	}
	for _, tc := range tt {
		keyValues, err := readFile(tc.fileName)
		if err != nil && tc.expected {
			t.Fatalf("unable to get to read the file")
		}
		if _, found := keyValues[tc.key]; !found && tc.expected {
			t.Errorf("could not find key, %s in map", tc.key)
		} else if found && !tc.expected {
			t.Errorf("found key %s in map, was not expecting it", tc.key)
		} else if found && tc.expected {
			if keyValues[tc.key] != tc.value {
				t.Errorf("found wrong value in map, expected: %s, got: %s", tc.value, keyValues[tc.key])
			}
		}
	}
}

func TestDetectLineContents(t *testing.T) {
	tt := []struct {
		name       string
		pattern    string
		line       string
		isComment  bool
		regExpFail bool
	}{
		{"Detect key value", `([a-zA-Z0-9_ ]+)=(.+)`, "this = that", true, false},
		{"Detect comment", `(\/\/).+`, "// Im a comment", true, false},
		{"Incorrect regex pattern", ")(", "", false, true},
	}

	for _, tc := range tt {
		lineIsComment, err := detectLineContents(tc.pattern, tc.line)
		if err != nil {
			if !tc.regExpFail {
				t.Errorf("unable to check if line %s was a comment or not %v", tc.line, err)
			}
		}
		if lineIsComment != tc.isComment {
			t.Errorf("%s - expected to see line was comment: %v, got: %v", tc.name, tc.isComment, lineIsComment)
		}
	}
}

func TestUnderscoreUppercase(t *testing.T) {
	tt := []struct {
		name, raw, expected string
		expectedParse       bool
	}{
		{"Single word", "word", "WORD", true},
		{"Single word whitespace", "    word    ", "WORD", true},
		{"Multiple word", "word of life", "WORD_OF_LIFE", true},
		{"Multiple word whitespace", "    word of life    ", "WORD_OF_LIFE", true},
		{"Mismatched entry", "something", "SOMETHING_ELSE", false},
	}

	for _, tc := range tt {
		processedString := underscoreUpperCase(tc.raw)
		if processedString != tc.expected {
			if tc.expectedParse {
				t.Errorf("%s - failed to process string, expected: %s, got: %s", tc.name, tc.expected, processedString)
			}
		}
	}
}

func TestGetPrefix(t *testing.T) {
	tt := []struct {
		prefix      string
		expectedGet bool
	}{
		{"", true},
		{"something", true},
	}

	defaultPrefix := "BL_SENV_PACKAGE"

	for _, tc := range tt {
		if tc.prefix != "" {
			err := setPrefix(tc.prefix)
			if err != nil {
				t.Errorf("error setting the prefix %v", err)
			}
		}
		prefix, err := getPrefix()
		if err != nil && !tc.expectedGet {
			t.Errorf("was unable to get the prefix var")
		}
		if tc.prefix == "" && prefix != defaultPrefix {
			t.Errorf("failed to get the correct env prefix, wanted: %s, got: %s", defaultPrefix, prefix)
		}
	}
}

func TestSetPrefix(t *testing.T) {
	tt := []struct {
		prefix      string
		expectedSet bool
	}{
		{"name", true},
		{"some name with spaces", true},
	}

	for _, tc := range tt {
		err := setPrefix(tc.prefix)
		if err != nil {
			if tc.expectedSet {
				t.Fatalf("unable to set prefix: %s, %v", tc.prefix, err)
			}
		}
		prefix, present := os.LookupEnv("SENV_PREFIX")
		if !present {
			t.Fatalf("failed to set env var %s", tc.prefix)
		} else if prefix != underscoreUpperCase(tc.prefix) {
			t.Errorf("prefix did not match wanted: %s, got: %v", tc.prefix, prefix)
		}
	}
}

func TestExtractKeyValue(t *testing.T) {
	tt := []struct {
		pattern      string
		line         string
		key          string
		value        string
		expectedPass bool
	}{
		{"([a-zA-Z0-9_ ]+)=(.+)", "simple = pass", "simple", "pass", true},
		{"([a-zA-Z0-9_ ]+)=(.+)", "nospace=pass", "nospace", "pass", true},
		{"([a-zA-Z0-9_ ]+)=(.+)", "nospace=pass", "nospace", "pass", true},
		{"([a-zA-Z0-9_ ]+)=(.+)", "    extraspace = pass     ", "extraspace", "pass", true},
		{"([a-zA-Z0-9_ ]+)=(.+)", "something else = pass", "set to fail", "some value", false},
		{"([a-zA-Z0-9_ ]+)=(.+)", "no value =", "no value", "some value", false},
		{"([a-zA-Z0-9_ ]+)=(.+)", "space value =  ", "space value", "some value", false},
		{"([a-zA-Z0-9_ ]+)=(.+)", "no valu", "no valu", "some value", false},
		{"([a-zA-Z0-9_ ]+)=(.+)", "=  no key", "no key", "some value", false},
		{"([a-zA-Z0-9_ ]+)=(.+)", "     =  space key", "space key", "some value", false},
		{`[`, "     =  invalid regexp", "space key", "some value", false},
	}
	for _, tc := range tt {
		key, value, err := extractKeyValue(tc.pattern, tc.line)
		if err != nil && tc.expectedPass {
			t.Errorf("unable to parse the line: %s, error: %v", tc.line, err)
		}
		if tc.key != key && tc.expectedPass || tc.value != value && tc.expectedPass {
			t.Errorf("key extraction failure expected, Key: %s, Value: %s, got - Key: %s, Value: %s", tc.key, tc.value, key, value)
		}
	}
}
