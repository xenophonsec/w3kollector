package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func saveLineToFile(filename string, content string) {
	if !dontsave {
		filePath := filepath.Join(outputDir, filename)
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Failed to open to file", filePath)
		} else {
			_, err := f.WriteString(content + "\n")
			if err != nil {
				fmt.Println("Failed to write to file", filePath)
			}
		}
		if err := f.Close(); err != nil {
			log.Fatal("Failed to close file: "+filePath, err)
		}
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
	if !dontsave {
		if outputFlag == "" {
			var err error
			outputDir, err = os.Getwd()
			if err != nil {
				panic(err)
			}
			outputDir = filepath.Join(outputDir, targetDomain)
			_, err = os.Stat(outputDir)
			if os.IsNotExist(err) {
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
}

// extract domain target domain from url
func urlToDomain(url string) string {
	url = strings.Replace(url, "https://", "", -1)
	targetDomain := strings.Replace(url, "http://", "", -1)
	slashIndex := strings.Index(targetDomain, "/")
	if slashIndex != -1 {
		targetDomain = targetDomain[0:slashIndex]
	}
	return targetDomain
}

func handleTargetURL(url string) string {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	return url
}
