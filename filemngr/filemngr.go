// Package filemngr provides functions for dealing with files
package filemngr

import (
	"apricate/log"
	"io/ioutil"
	"os"
	"strings"
)

// TODO: Test
// Delete file if it exists, else continue (ignore errors)
func DeleteIfExists(path string) {
	os.Remove(path)
}

// TODO: Test
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

// TODO: Test
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

// TODO: Test
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

// TODO: Test
// Search slice for search key, returns list of all lines where found
func GetLinesContainingKey(searchKey string, lines []string) ([]string) {
	res := []string{}
	for i, line := range lines {
		if strings.Contains(line, searchKey) {
			log.Debug.Printf("Found search key %s at line: %v", searchKey, i)
			res = append(res, line)
		}
	} 
	return res
}

// TODO: Test
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

// TODO: Test
// Converts os.File to bytes slice
func ConvertFileToBytes(file *os.File) []byte {
	log.Debug.Println("Reading from file")
	byteValue, _ := ioutil.ReadAll(file)
	return byteValue
}

// TODO: Test
// Writes file with specified bytes
func WriteBytesToFile(path string, bytes []byte) error {
	log.Debug.Printf("Writing bytes to file at %s", path)
	err := ioutil.WriteFile(path, bytes, 0644)
	log.Debug.Printf("Finished, err?: %v", err)
	return err
}

// TODO: Test
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

// TODO: Test
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

// Searches a slice of strings for a specified key-value pair and returns the value if it exists
func GetKeyFromLines(key string, lines []string) string {
	found := GetLinesContainingKey(key, lines)
	if len(found) < 1 {
		log.Info.Printf("Key not found: %s", key)
		return ""
	}
	res := []string{}
	for _, line := range found {
		kv_pair := strings.Split(line, "=")
		if len(kv_pair) < 2 {
			continue
		}
		if key != kv_pair[0] { // Additional check for case where key is substring or found later in the line
			continue
		}
		value := strings.Join(kv_pair[1:], "=") // join in-case value has an = for some reason
		res = append(res, value)
	}
	if len(res) < 1 {
		log.Info.Printf("Key not found: %s", key)
		return ""
	} else if len(res) > 1 {
		log.Error.Printf("Multiple matching keys (%s) in file, returning none: %v", key, res)
		return ""
	} else {
		return res[0]
	}
}