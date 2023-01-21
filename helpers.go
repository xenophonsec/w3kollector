package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
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

func handleOutputPath(outputFlag string, targetDomain string) {
	if outputFlag == "" {
		var err error
		outputDir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		outputDir = filepath.Join(outputDir, targetDomain)
		_, err = os.Stat(outputDir)
		if os.IsNotExist(err) {
			color.Red(outputDir + " does not exist")
			os.Mkdir(outputDir, os.FileMode(0644))
		} else if err != nil {
			panic(err)
		}
	} else {
		outputDir = outputFlag
		_, err := os.Stat(outputDir)
		if os.IsNotExist(err) {
			color.Red(outputDir + " does not exist")
			os.Exit(1)
		} else if err != nil {
			panic(err)
		}
	}
}
