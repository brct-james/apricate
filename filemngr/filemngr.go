package filemngr

import (
	"io/ioutil"
	"os"
	"strings"

	"apricate/log"
)

// Delete file if it exists, else continue
func DeleteIfExists(path string) {
	os.Remove(path)
}

// Ensure file exists, if not create it
func Touch(name string) error {
	log.Debug.Printf("Ensuring %s exists", name)
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		// Depending on the file different responses may be valid - pass errors up the stack
		return err
	}
	return file.Close()
}

// Reads file at string path to a slice of strings by line
func ReadFileToLineSlice(filePath string) ([]string, error) {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		// Depending on the file different responses may be valid - pass errors up the stack
		return nil, err
	}
	lines := strings.Split(string(input), "\n")
	return lines, nil
}

// Search slice for search key, returns true, index if found, else false
func KeyInSliceOfLines(searchKey string, lines []string) (bool, int) {
	for i, line := range lines {
		if strings.Contains(line, searchKey) {
			log.Debug.Printf("Found search key %s at line: %v", searchKey, i)
			return true, i
		}
	}
	return false, 0
}

// Write slice of lines to file at path
func WriteLinesToFile(filePath string, lines []string) error {
	output := strings.Join(lines, "\n")
	err := ioutil.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		// Depending on the file different responses may be valid - pass errors up the stack
		return err
	}
	return nil
}

// Converts os.File to bytes slice
func ConvertFileToBytes(file *os.File) []byte {
	log.Debug.Println("Reading from file")
	byteValue, _ := ioutil.ReadAll(file)
	return byteValue
}

// Writes file with specified bytes
func WriteBytesToFile(path string, bytes []byte) error {
	log.Debug.Printf("Writing bytes to file at %s", path)
	err := ioutil.WriteFile(path, bytes, 0644)
	log.Debug.Printf("Finished, err?: %v", err)
	return err
}

// Reads specified file, returns byte slice
func ReadFileToBytes(path string) ([]byte, error) {
	readFile, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	log.Debug.Println("Successfully opened file for reading to bytes at: " + path)
	defer readFile.Close()
	return ConvertFileToBytes(readFile), nil
}

// Reads every file in directory, returning slice of bytevalues
func ReadFilesToBytes(path_to_directory string) ([][]byte, error) {
	files, err := ioutil.ReadDir(path_to_directory)
	if err != nil {
		return [][]byte{}, err
	}
	bytes := make([][]byte, len(files))
	for i, file := range files {
		filename := file.Name()
		var readErr error
		bytes[i], readErr = ReadFileToBytes(path_to_directory + "/" + filename)
		if readErr != nil {
			return [][]byte{}, readErr
		}
	}
	return bytes, nil
}