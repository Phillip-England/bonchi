package bonchi

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

func BundleJs(inputDir string, out string) (string, error) {
	var files []string
	err := filepath.Walk(inputDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".js") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Slice(files, func(i, j int) bool {
		numI := extractPrefixNumber(files[i])
		numJ := extractPrefixNumber(files[j])
		return numI < numJ
	})
	output := ""
	for _, file := range files {
		jsBytes, err := os.ReadFile(file)
		if err != nil {
			return output, err
		}
		output += string(jsBytes) + "\n"
	}

	m := prepareMinify()
	minifiedJS, err := m.String("application/javascript", output)
	if err != nil {
		return "", fmt.Errorf("failed to minify JavaScript: %w", err)
	}
	err = writeToFile(out, minifiedJS, true)
	if err != nil {
		return minifiedJS, err
	}
	return minifiedJS, nil
}

func BundleCss(inputDir string, out string) (string, error) {
	var files []string
	err := filepath.Walk(inputDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".css") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Slice(files, func(i, j int) bool {
		numI := extractPrefixNumber(files[i])
		numJ := extractPrefixNumber(files[j])
		return numI < numJ
	})
	output := ""
	for _, file := range files {
		cssBytes, err := os.ReadFile(file)
		if err != nil {
			return output, err
		}
		output += string(cssBytes) + "\n"
	}
	output = handleBonchiMix(output)
	m := prepareMinify()
	minifiedCSS, err := m.String("text/css", output)
	if err != nil {
		return "", fmt.Errorf("failed to minify CSS: %w", err)
	}
	err = writeToFile(out, minifiedCSS, true)
	if err != nil {
		return minifiedCSS, err
	}
	return minifiedCSS, nil
}

func extractPrefixNumber(filePath string) int {
	base := filepath.Base(filePath)
	parts := strings.Split(base, ".")
	if len(parts) < 2 {
		return 999999
	}
	num, err := strconv.Atoi(parts[0])
	if err != nil {
		return 999999
	}
	return num
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

func prepareMinify() *minify.M {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	return m
}

func minifyStaticFiles(m *minify.M, dirPath string) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		var mimetype string
		switch ext {
		case ".css":
			mimetype = "text/css"
		case ".html":
			mimetype = "text/html"
		case ".js":
			mimetype = "application/javascript"
		case ".json":
			mimetype = "application/json"
		case ".svg":
			mimetype = "image/svg+xml"
		case ".xml":
			mimetype = "text/xml"
		default:
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err)
			return err
		}
		defer f.Close()
		fileBytes, err := io.ReadAll(f)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			return err
		}
		minifiedBytes, err := m.Bytes(mimetype, fileBytes)
		if err != nil {
			fmt.Printf("Error minifying file: %s\n", err)
			return err
		}
		err = os.WriteFile(path, minifiedBytes, info.Mode()) // Preserving original file permissions
		if err != nil {
			fmt.Printf("Error writing minified file: %s\n", err)
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the directory: %s\n", err)
	}
}
