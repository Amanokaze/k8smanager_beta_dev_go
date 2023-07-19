package logger

import "os"

// CreateDir is a function that creates a directory.
// path : directory path
func CreateDir(path string, user string) (bool, error) {
	if !IsExistFile(path) {
		err := os.Mkdir(path, 0755)
		return CheckError(err), err
	}

	return true, nil
}
