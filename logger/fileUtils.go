package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// CheckError is a function that checks error.
func CheckError(err error) bool {
	if err != nil {
		log.Println(err.Error())
		return true
	}

	return false
}

// CreateFile is a function that creates a file.
// If the file already exists, it will be truncated.
// path : file path
func CreateFile(path string) error {
	// check if file exists
	var _, err = os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if CheckError(err) {
			return err
		}
		defer file.Close()
	}

	return nil
}

// IsExistFile is a function that checks if a file exists.
// fname : file name
func IsExistFile(fname string) bool {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	}

	return true
}

// CheckOrCreateFilePathDir is a function that checks or creates a directory.
// filePath : directory path
func CheckOrCreateFilePathDir(filePath string, user string) error {
	tmp := strings.Replace(filePath, "\\", "/", -1)
	//tmp := filePath
	nPos := 0
	for {
		pos := strings.Index(tmp, "/")
		if pos > 0 {
			dir := filePath[:pos+nPos]
			_, err := CreateDir(dir, user)
			if err != nil {
				return err
			}

			tmp = tmp[pos+1:]
			nPos += pos + 1
		} else if pos == 0 {
			tmp = tmp[pos+1:]
			nPos += pos + 1
		} else {
			break
		}
	}

	return nil
}

// DeleteFile is a function that deletes a file.
// path : file path
func DeleteFile(path string) {
	// delete file
	if IsExistFile(path) {
		var err = os.Remove(path)
		CheckError(err)
	}
}

// RenameFile is a function that renames a file.
// orgPath : original file path
// newPath : new file path
func RenameFile(orgPath string, newPath string) bool {
	err := os.Rename(orgPath, newPath)
	return CheckError(err)
}

// ReadFile is a function that reads a file.
// path : file path
func ReadFile(path string) (string, error) {
	// Open file for reading.
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read file, line by line
	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)
		// Break if finally arrived at end of file
		if err == io.EOF {
			break
		}
		// Break if error occured
		if err != nil && err != io.EOF {
			if err != nil {
				return "", err
			}

			break
		}
	}
	// fmt.Println("Reading from file.")
	// fmt.Println(string(text))

	return string(text), nil
}

// WriteBytes is a function that writes bytes to a file.
// path : file path
// bytes : bytes to write
func WriteBytes(path string, bytes []byte) {
	file, err := os.Create(path) // hello.txt 파일 생성
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close() // main 함수가 끝나기 직전에 파일을 닫음

	_, err = file.Write(bytes) // s를 []byte 바이트 슬라이스로 변환, s를 파일에 저장
	if err != nil {
		fmt.Println(err)
		return
	}
}

// WriteFile is a function that writes a file.
// path : file path
// text : text to write
func WriteFile(path string, text string) error {
	// Open file using READ & WRITE permission.
	//var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	var file, err = os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return err
	}

	// Save file changes.
	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

// AppendFile is a function that appends a file.
// path : file path
// text : text to append
func AppendFile(path string, text string) error {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return err
	}

	// Save file changes.
	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

// ReadFileLastBytes is a function that reads last bytes from a file.
// path : file path
// lastlen : last bytes length
func ReadFileLastBytes(path string, lastlen int64) []byte {
	file, err := os.OpenFile(
		path,
		os.O_RDONLY,
		os.FileMode(0644)) // 파일 권한은 644
	if CheckError(err) {
		return nil
	}
	defer file.Close() // main 함수가 끝나기 직전에 파일을 닫음

	fi, err := file.Stat() // 파일 정보 가져오기
	if CheckError(err) {
		return nil
	}

	var data = make([]byte, lastlen) // 파일 크기만큼 바이트 슬라이스 생성
	_, err = file.Seek(fi.Size()-lastlen, 5)
	if CheckError(err) {
		return nil
	}

	_, err = file.Read(data)
	if CheckError(err) {
		return nil
	}

	return data
}

// WriteBufFileLog is a function that writes a file.
// buf : bytes buffer
// filepath : file path
func WriteBufFileLog(buf *bytes.Buffer, filepath string) error {
	var file, err = os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if CheckError(err) {
		fmt.Printf("Logfile Open Error - %s %v", filepath, err.Error())
		return err
	}

	defer file.Close()

	_, err = file.WriteString(buf.String())
	if CheckError(err) {
		fmt.Printf("Logfile Write Error - %s %v", filepath, err.Error())
		return err
	}
	//Save file changes.
	err = file.Sync()
	if CheckError(err) {
		fmt.Printf("Logfile Save Change Error - %s %v", filepath, err.Error())
		return err
	}

	return nil
}
