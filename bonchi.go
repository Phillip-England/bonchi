package bonchi

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Bundle(inputDir string, out string) (string, error) {
	output := ""
	err := filepath.Walk(inputDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		pathParts := strings.Split(path, ".")
		lastPart := pathParts[len(pathParts)-1]
		if lastPart == "css" {
			cssBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			cssStr := string(cssBytes)
			output += cssStr + "\n"
		}
		return nil
	})
	if err != nil {
		return output, err
	}
	output = handleBonchiMix(output)
	err = writeToFile(out, output, true)
	if err != nil {
		return output, err
	}
	return output, nil
}

func getBundledCssFromFile(inputPath string) (string, error) {
	output := ""
	fileBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return "", err
	}
	fileStr := string(fileBytes)
	lines := strings.Split(fileStr, "\n")
	for _, line := range lines {
		lineParts := strings.Split(line, " ")
		if len(lineParts) != 2 {
			output += line + "\n"
			continue
		}
		if lineParts[0] != "@bonchi" {
			output += line + "\n"
			continue
		}
		filePath := strings.ReplaceAll(lineParts[1], ";", "")
		cssBytes, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		cssContent := string(cssBytes)
		output += cssContent + "\n"
	}
	return output, nil
}

func handleBonchiMix(css string) string {
	out := ""
	lines := strings.Split(css, "\n")
	for _, line := range lines {
		if strings.Contains(line, "bonchi-mix") {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}
			classes := parts[1]
			classes = strings.ReplaceAll(classes, ";", "")
			classes = strings.ReplaceAll(classes, "'", "")
			classes = strings.ReplaceAll(classes, "\"", "")
			classNames := strings.Split(classes, " ")
			for _, name := range classNames {
				classCss := GetClassCssByName(css, name)
				out += classCss
			}
			continue
		}
		out += line + "\n"
	}
	return out
}

func GetClassCssByName(css string, className string) string {
	classContent := ""
	isInsideClass := false
	lines := strings.Split(css, "\n")
	for _, line := range lines {
		sqLine := strings.ReplaceAll(line, " ", "")
		locator := className + "{"
		if strings.Contains(sqLine, "}") {
			isInsideClass = false
		}
		if strings.Contains(sqLine, locator) {
			isInsideClass = true
			continue
		}
		if isInsideClass {
			classContent += line + "\n"
		}
	}
	return classContent
}

func writeToFile(path string, content string, overwrite bool) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	if _, err := os.Stat(path); err == nil {
		if !overwrite {
			return errors.New("file already exists and overwrite is false")
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check file existence: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
