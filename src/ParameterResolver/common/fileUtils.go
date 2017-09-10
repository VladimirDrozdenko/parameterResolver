package common

import (
	"os"
	"errors"
	"io/ioutil"
)

// checks if file is less than MaxFileSizeInBytes and returns error if it is not
func ValidateFileAndSize(source string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	fileStats, err := file.Stat()
	if err != nil {
		return err
	}
	if fileStats.Size() > MaxFileSizeInBytes {
		return errors.New("File is too large.")
	}
	return nil
}

// returns the text inside a given (relative ?) path to a file
func ReadTextFromFile(source string) (string, error) {
	dat, err := ioutil.ReadFile(source)
	if err != nil {
		return "", err
	}

	unresolvedText := string(dat)
	return unresolvedText, nil
}

func WriteToFile(resolvedText string, destination string) error {
	f, err := os.Create(destination)

	defer f.Close()

	if err != nil {
		return err
	}
	_, err = f.WriteString(resolvedText)
	if err != nil {
		return err
	}

	return nil
}
