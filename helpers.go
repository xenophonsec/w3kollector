package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func saveLineToFile(filename string, content string) {
	filePath := filepath.Join(outputDir, filename)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to write to file", filename)
	} else {
		_, err := f.WriteString(content + "\n")
		if err != nil {
			fmt.Println("Failed to write to file", filename)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal("error", err)
	}
}

func arrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
