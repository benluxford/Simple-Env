package senv

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Load : Loads all the environment variables from the file.
func Load(filePath string) (err error) {
	// get all the env variables
	keyValues, err := readFile(filePath)
	// if the prefix has been defined - set it first
	if _, found := keyValues["SENV_PREFIX"]; found {
		err = setPrefix(keyValues["SENV_PREFIX"])
	}
	// loop over entire map and set all values
	for k, v := range keyValues {
		err = SetVar(k, v)
	}
	return
}

// SetVar : Sets the environment variable or returns error
func SetVar(key, value string) (err error) {
	// Find and get the prefix from the env
	prefix, err := getPrefix()
	// add the prefix to the key
	key = fmt.Sprintf("%s_%s", prefix, key)
	// format the key
	key = underscoreUpperCase(key)
	// set the key value
	err = os.Setenv(key, value)
	return
}

// GetVar : Returns environment variable or error
func GetVar(key string) (value string, err error) {
	// Find and get the prefix from the env
	prefix, err := getPrefix()
	// add the prefix to the key
	key = fmt.Sprintf("%s_%s", prefix, key)
	// format the key
	key = underscoreUpperCase(key)
	// get the value
	value, present := os.LookupEnv(key)
	// if no value, set the error
	if !present {
		err = fmt.Errorf("%s was not set in env vars", value)
	}
	return
}

// setPrefix : sets the prefix value in the env, returns error/nil
func setPrefix(prefix string) (err error) {
	// set the prefix
	err = os.Setenv("SENV_PREFIX", underscoreUpperCase(prefix))
	return
}

// getPrefix : gets or creates the default prefix from the env variables
func getPrefix() (prefix string, err error) {
	defaultPrefix := "BL_SENV_PACKAGE"
	// Find and get the prefix from the env
	prefix, present := os.LookupEnv("SENV_PREFIX")
	// If the prefix is missing throw error
	if !present {
		// If not present create it
		err = setPrefix(defaultPrefix)
		// set the return value to the default prefix
		prefix = defaultPrefix
	}
	return
}

// detectLineContents : detects if the current line is a comment or not
func detectLineContents(pattern, line string) (detected bool, err error) {
	// get the regex ready
	commentTest, err := regexp.Compile(pattern)
	if err != nil {
		err = fmt.Errorf("Line detection Package ENV: %v", err)
		return
	}
	// return true if the line matched the pattern
	detected = commentTest.MatchString(line)
	return
}

// underscoreUpperCase : takes a given string and parses it as UPPERCASE_WITH_UNDERSCORES
func underscoreUpperCase(raw string) string {
	// Trim any whitespace
	trimmed := strings.TrimSpace(raw)
	// Uppercase the string
	upperCase := strings.ToUpper(trimmed)
	// Return string with all spaces replaced with underscores
	return strings.Replace(upperCase, " ", "_", -1)
}

// extractKeyValue : returns the trimmed key and value pair from a given line
func extractKeyValue(pattern, line string) (key, value string, err error) {
	// compile the regexp to test against
	keyValueTest, err := regexp.Compile(pattern)
	if err != nil {
		err = fmt.Errorf("Key Value extraction ENV: %v", err)
		return
	}
	// test pattern matches string (already tested)
	if keyValueTest.MatchString(line) {
		// find all submatches
		results := keyValueTest.FindStringSubmatch(line)
		// set values and trim whitespace
		key = strings.TrimSpace(results[1])
		value = strings.TrimSpace(results[2])
		return
	}
	err = fmt.Errorf("unable to match string to pattern")
	return
}

// readFile : reads the env file from its location and parses the contents into a k,v map
func readFile(AbsoluteLocation string) (keyValues map[string]string, err error) {
	// create empty map
	keyValues = make(map[string]string)
	// open the file
	file, err := os.Open(AbsoluteLocation)
	if err != nil {
		return
	}
	// create new scanner
	scanner := bufio.NewScanner(file)
	// for line in file
	for scanner.Scan() {
		// extract the text
		text := scanner.Text()
		// test, text is key value pair
		if isKeyValue, err := detectLineContents("([a-zA-Z0-9_ ]+)=(.+)", text); err != nil {
			break
		} else if err == nil && isKeyValue {
			// extrace the k,v pair from the string
			key, value, err := extractKeyValue("([a-zA-Z0-9_ ]+)=(.+)", text)
			if err != nil {
				break
			}
			// test that value was not predefined in document
			if _, found := keyValues[key]; !found {
				// add value to k,v map
				keyValues[key] = value
			} else {
				err = fmt.Errorf("key value [%s: %s], already exists, duplicate entry in file", key, value)
				break
			}
		}
	}
	return
}
