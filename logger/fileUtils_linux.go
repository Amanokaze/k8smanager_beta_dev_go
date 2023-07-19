package logger

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// CreateDir is a function that creates a directory.
// path : directory path
func CreateDir(path string, os_user string) (bool, error) {
	if !IsExistFile(path) {
		err := os.Mkdir(path, 0755)
		CheckError(err)

		if os_user != "" {
			// Change the owner of the directory to the user who runs the program.
			group, err := user.Lookup(os_user)
			if err != nil {
				return CheckError(err), err
			}

			uid, err := strconv.Atoi(group.Uid)
			if err != nil {
				return CheckError(err), err
			}

			gid, err := strconv.Atoi(group.Gid)
			if err != nil {
				return CheckError(err), err
			}

			err = syscall.Chown(path, uid, gid)
			if err != nil {
				return CheckError(err), err
			}
		}
	}

	return true, nil
}
