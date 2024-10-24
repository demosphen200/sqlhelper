package ex

import (
	"bufio"
	"os"
	"path/filepath"
)

func FileNameWithoutExt(entry os.DirEntry) string {
	var filename = entry.Name()
	var ext = filepath.Ext(filename)
	return filename[0 : len(filename)-len(ext)]
}

func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func IsFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}
