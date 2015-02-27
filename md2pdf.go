package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	HELP = `md2pdf [-h] [-x] [-s] [-i dir] [-o file] file.md
Transform a given Markdown file into PDF.
-h        To print this help page.
-x        Print intermediate XHTML output.
-s        Print stylesheet used for transformation.
-a        Output article (instead of blog entry).
-i dir    To indicate image directory.
-o file   The name of the file to output.
file.md   The markdown file to convert.
Note:
This program calls pandoc, xsltproc and htmldoc that must have been installed.`
	STYLESHEET = `<?xml version="1.0" encoding="utf-8"?>

<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                version="1.0">

  <xsl:output method="xml" encoding="ISO-8859-1"/>
  <xsl:param name="id">ID</xsl:param>
  <xsl:param name="date">DATE</xsl:param>
  <xsl:param name="title">TITLE</xsl:param>
  <xsl:param name="author">AUTHOR</xsl:param>
  <xsl:param name="email">EMAIL</xsl:param>
  <xsl:param name="lang">fr</xsl:param>
  <xsl:param name="toc">yes</xsl:param>

  <!-- catch the root element -->
  <xsl:template match="/xhtml">
    <xsl:apply-templates/>
  </xsl:template>

  <xsl:template match="@*|node()">
    <xsl:copy>
      <xsl:apply-templates select="@*|node()"/>
    </xsl:copy>
  </xsl:template>

  <xsl:template match="table">
    <table border='1' cellpadding='5'>
      <xsl:apply-templates/>
	</table>
	<p></p>
  </xsl:template>

  <xsl:template match="pre">
    <table width='100%' border='0' cellpadding='10'>
	  <tr>
	    <td bgcolor='#F0F0F0'>
		  <pre>
		    <xsl:apply-templates/>
		  </pre>
		</td>
	  </tr>
	</table>
	<p></p>
  </xsl:template>

</xsl:stylesheet>`
	XHTML_HEADER = "<xhtml>\n<body>\n"
	XHTML_FOOTER = "</body>\n</xhtml>"
)

func processXsl(tmpFile string, data map[string]string) {
	xslFile, err := ioutil.TempFile("/tmp", "md2pdf-")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(xslFile.Name(), []byte(STYLESHEET), 0644)
	if err != nil {
		panic(err)
	}
	defer os.Remove(xslFile.Name())
	params := make([]string, 0, 2+3*len(data))
	for name, value := range data {
		params = append(params, "--stringparam")
		params = append(params, name)
		params = append(params, value)
	}
	params = append(params, xslFile.Name())
	params = append(params, tmpFile)
	command := exec.Command("xsltproc", params...)
	result, err := command.CombinedOutput()
	if err != nil {
		println(result)
		panic(err)
	}
	err = ioutil.WriteFile(tmpFile, result, 0644)
	if err != nil {
		panic(err)
	}
}

func markdown2xhtml(markdown string) []byte {
	mdFile, err := ioutil.TempFile("/tmp", "md2xsl-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(mdFile.Name())
	ioutil.WriteFile(mdFile.Name(), []byte(markdown), 0644)
	command := exec.Command("pandoc", mdFile.Name(), "-f", "markdown", "-t", "html")
	result, err := command.CombinedOutput()
	if err != nil {
		println(result)
		panic(err)
	}
	return []byte(XHTML_HEADER + string(result) + XHTML_FOOTER)
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

func imageDir(text, imgDir string) string {
	r := regexp.MustCompile(`!\[(.*?)\]\((.*?/)*(.*?)\)`)
	if len(imgDir) > 0 {
		return r.ReplaceAllString(text, "![$1]("+imgDir+"/$3)")
	} else {
		return r.ReplaceAllString(text, "![$1]($3)")
	}
}

func generatePdf(xhtmlFile, outFile string) {
	params := []string{
		"--outfile", outFile,
		"--size", "A4",
		"--top", "2cm",
		"--bottom", "2cm",
		"--left", "2cm",
		"--right", "2cm",
		"--bodyfont", "Times",
		"--fontsize", "12",
		"--header", "...",
		"--footer", "dt1",
		"--headfootfont", "Courier-Oblique",
		"--headfootsize", "10",
		"--linkcolor", "#0000A0",
		"--linkstyle", "plain",
		"--permissions", "no-modify",
		"--charset", "iso-8859-1",
		"--no-title",
		"--no-toc",
		"--compression=9",
		"--embedfonts",
		"--webpage",
		xhtmlFile,
	}
	command := exec.Command("htmldoc", params...)
	result, err := command.CombinedOutput()
	if err != nil {
		println(string(result))
		panic(err)
	}
}

func processFile(filename string, printXhtml bool, imgDir, outFile string) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	data, markdown := markdownData(string(source))
	markdown = imageDir(markdown, imgDir)
	xhtml := markdown2xhtml(markdown)
	if printXhtml {
		fmt.Println(string(xhtml))
		return
	}
	tmpFile, err := ioutil.TempFile("/tmp", "md2pdf-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())
	ioutil.WriteFile(tmpFile.Name(), xhtml, 0644)
	processXsl(tmpFile.Name(), data)
	if len(outFile) == 0 {
		outFile = filename[0:len(filename)-len(filepath.Ext(filename))] + ".pdf"
	}
	generatePdf(tmpFile.Name(), outFile)
}

func main() {
	xhtml := false
	imgDir := ""
	outFile := ""
	file := ""
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
		} else if arg == "-x" || arg == "--xhtml" {
			xhtml = true
		} else if arg == "-s" || arg == "--stylesheet" {
			fmt.Println(STYLESHEET)
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
	processFile(file, xhtml, imgDir, outFile)
}
