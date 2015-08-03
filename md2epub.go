package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	HELP = `md2epub [-h] [-i dir] [-o file] file.md
Transform a given Markdown file into XML.
-h        To print this help page.
-i dir    To indicate image directory.
-o file   The name of the file to output.
file.md   The markdown file to convert.
Note: this program calls pandoc that must have been installed.`
)

func printError(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func markdownData(text string) (map[string]string, string) {
	data := make(map[string]string)
	lines := strings.Split(text, "\n")
	var limit int
	for index, line := range lines {
		if strings.HasPrefix(line, "% ") && strings.Index(line, ":") >= 0 {
			name := strings.TrimSpace(line[2:strings.Index(line, ":")])
			value := strings.TrimSpace(line[strings.Index(line, ":")+1 : len(line)])
			data[name] = value
		} else {
			limit = index
			break
		}
	}
	return data, strings.Join(lines[limit:len(lines)], "\n")
}

func insertMeta(text string, meta map[string]string) string {
	header := ""
	if title, ok := meta["title"]; ok {
		header += title + "\n"
	} else {
		printError("Missing title in metadata")
	}
	if author, ok := meta["author"]; ok {
		header += author + "\n"
	} else {
		printError("Missing author in metadata")
	}
	return header + text
}

func imageDir(text, imgDir string) string {
	r := regexp.MustCompile(`!\[(.*?)\]\((.*?/)*(.*?)\)`)
	if len(imgDir) > 0 {
		return r.ReplaceAllString(text, "![$1]("+imgDir+"/$3)")
	} else {
		return r.ReplaceAllString(text, "![$1]($3)")
	}
}

func markdown2epub(markdown, outFile string) {
	mdFile, err := ioutil.TempFile("/tmp", "md2epub-")
	if err != nil {
		printError("Error creating temporary file")
	}
	defer os.Remove(mdFile.Name())
	err = ioutil.WriteFile(mdFile.Name(), []byte(markdown), 0644)
	if err != nil {
		printError("Error writing temporary file")
	}
	command := exec.Command("pandoc", "-f", "markdown", "-t", "epub",
		"-o", outFile, mdFile.Name())
	result, err := command.CombinedOutput()
	if err != nil {
		output := strings.TrimSpace(string(result))
		printError(fmt.Sprintf("Error calling pandoc:\n%s", output))
	}
}

func processFile(filename, imgDir, outFile string) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		printError(fmt.Sprintf("Error reading file %s", filename))
	}
	data, markdown := markdownData(string(source))
	markdown = imageDir(markdown, imgDir)
	markdown = insertMeta(markdown, data)
	if outFile == "" {
		outFile = filename[:strings.LastIndex(filename, ".")] + ".epub"
	}
	markdown2epub(markdown, outFile)
}

func main() {
	file := ""
	imgDir := ""
	outFile := ""
	if len(os.Args) < 2 {
		fmt.Println(HELP)
		os.Exit(1)
	}
	skip := false
	args := os.Args[1:]
	for i, arg := range args {
		if skip {
			skip = false
			continue
		}
		if arg == "-h" || arg == "--help" {
			fmt.Println(HELP)
			os.Exit(0)
		} else if arg == "-i" || arg == "--image-dir" {
			imgDir = args[i+1]
			skip = true
		} else if arg == "-o" || arg == "--output-file" {
			outFile = args[i+1]
			skip = true
		} else {
			file = arg
		}
	}
	processFile(file, imgDir, outFile)
}
