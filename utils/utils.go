package utils

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
)

// line feed
var LineFeed = []byte{10}

// DSS file name to be excluded
const DSSFileName = ".DS_Store"

// GetContents gets the contents separated by the line feed from the given file path.
// Empty lines ignored
func GetContents(path string) (contents [][]byte, err error) {
	bz, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	for _, line := range bytes.Split(bz, LineFeed) {
		if len(line) != 0 {
			contents = append(contents, line)
		}
	}

	return contents, nil
}

// GetFiles gets files from the given path.
// Empty files ingored and subdirectories included
func GetFiles(path string) (files [][]byte, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files = make([][]byte, 0)

	for _, entry := range entries {
		if !entry.IsDir() {
			if entry.Name() != DSSFileName {
				file, err := os.ReadFile(fmt.Sprintf("%s/%s", path, entry.Name()))
				if err != nil {
					return nil, err
				}

				if len(file) != 0 {
					files = append(files, file)
				}
			}
		} else {
			subDirFiles, err := GetFiles(fmt.Sprintf("%s/%s", path, entry.Name()))
			if err != nil {
				return nil, err
			}

			if len(subDirFiles) != 0 {
				files = append(files, subDirFiles...)
			}
		}
	}

	return files, nil
}

// GetSubFilePaths gets the file paths located in the given directory
func GetSubFilePaths(path string) (paths []string, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	paths = make([]string, len(entries))

	for _, entry := range entries {
		paths = append(paths, entry.Name())
	}

	return paths, nil
}

// WriteFile writes the specified contents to the given file
func WriteFile(path string, contents []string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, content := range contents {
		f.WriteString(fmt.Sprintf("%s\n", content))
	}

	f.WriteString("\n")

	return nil
}

// SHA256 returns the sha256 hash of the given data
func SHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
