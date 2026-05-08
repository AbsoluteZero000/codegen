package main

import (
	"os"
	"path/filepath"
)

func readFile(filePath string) string {
	wd, err := os.Getwd()
	check(err)

	path := filepath.Join(wd, filePath)
	data, err := os.ReadFile(path)
	check(err)

	return string(data)
}
