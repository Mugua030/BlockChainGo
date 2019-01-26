package main

import "os"

func isFileExist(fileName string) bool {

	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
