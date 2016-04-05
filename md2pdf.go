package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	HELP = `md2pdf [-h] [-x] [-s] [-t] [-i dir] [-o file] file.md
Transform a given Markdown file into PDF.
-h        To print this help page.
-x        Print intermediate XHTML output.
-s        Print stylesheet used for transformation.
-t        Print html output.
-i dir    To indicate image directory.
-o file   The name of the file to output.
file.md   The markdown file to convert.
Note:
This program calls pandoc, xsltproc, htmldoc and faketime that must have been
installed.`
	STYLESHEET = `<?xml version="1.0" encoding="utf-8"?>

<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                version="1.0">

  <xsl:output method="xml" encoding="ISO-8859-1"/>
  <xsl:param name="title"/>
  <xsl:param name="author"/>
  <xsl:param name="date"/>
  <xsl:param name="id"/>
  <xsl:param name="email"/>
  <xsl:param name="lang"/>
  <xsl:param name="toc"/>

  <!-- catch the root element -->
  <xsl:template match="/xhtml/body">
    <head>
	  <meta http-equiv='Content-Type' content='text/xhtml; charset=ISO-8859-1'/>
	  <xsl:if test="$title">
	    <title><xsl:value-of select="$title"/></title>
	  </xsl:if>
	</head>
	<body>
    <xsl:if test="$title">
	  <center>
	    <h1><xsl:value-of select="$title"/></h1>
	  </center>
	  <p align="center"><i>
	    <font size="-1">
        <xsl:if test="$author">
	      <xsl:value-of select="$author"/><br/>
        </xsl:if>
	    <xsl:if test="$email">
		  <a href="mailto:{$email}">
		    <xsl:value-of select="$email"/>
		  </a>
	    </xsl:if>
		</font>
	  </i></p>
	  <br/>
	</xsl:if>
    <xsl:apply-templates/>
	</body>
  </xsl:template>

  <xsl:template match="@*|node()">
    <xsl:copy>
      <xsl:apply-templates select="@*|node()"/>
    </xsl:copy>
  </xsl:template>

  <xsl:template match="table">
	<center>
      <table border='1' cellpadding='5'>
        <xsl:apply-templates/>
	  </table>
	</center>
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

  <xsl:template match="p[img]">
    <p align="center" width="80%">
      <xsl:apply-templates />
    </p>
  </xsl:template>

  <xsl:template match="img">
    <img src="{@src}"><br/><i><xsl:value-of select="@alt"/></i></img>
  </xsl:template>

</xsl:stylesheet>`
	XHTML_HEADER = "<xhtml>\n<body>\n"
	XHTML_FOOTER = "</body>\n</xhtml>"
)

type MetaData struct {
	Title  string
	Author string
	Date   string
	Tags   []string
	Id     string
	Email  string
	Lang   string
	Toc    string
}

func printError(err error, message string) {
	if err != nil {
		println(message)
		os.Exit(1)
	}
}

func executeCommand(command *exec.Cmd) []byte {
	result, err := command.CombinedOutput()
	if err != nil {
		println(fmt.Sprintf("Error running command \"%s\":", strings.Join(command.Args, " ")))
		println(string(result))
		os.Exit(1)
		return []byte{}
	} else {
		return result
	}
}

func (d MetaData) ToMap() map[string]string {
	data := make(map[string]string)
	if d.Title != "" {
		data["title"] = d.Title
	}
	if d.Author != "" {
		data["author"] = d.Author
	}
	if d.Date != "" {
		data["date"] = d.Date
	}
	if len(d.Tags) != 0 {
		data["tags"] = strings.Join(d.Tags, ", ")
	}
	if d.Id != "" {
		data["id"] = d.Id
	}
	if d.Email != "" {
		data["email"] = d.Email
	}
	if d.Lang != "" {
		data["lang"] = d.Lang
	}
	if d.Toc != "" {
		data["toc"] = d.Toc
	}
	return data
}

var LOCALE = map[string]string{
	"fr": "fr_FR.UTF-8",
	"en": "en_US.UTF-8",
}

func processXsl(tmpFile string, data map[string]string) {
	xslFile, err := ioutil.TempFile("/tmp", "md2pdf-")
	printError(err, "Error creating XSL temporary file")
	err = ioutil.WriteFile(xslFile.Name(), []byte(STYLESHEET), 0644)
	printError(err, "Error writing XSL temporary file")
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
	result := executeCommand(command)
	err = ioutil.WriteFile(tmpFile, result, 0644)
	printError(err, "Error writing temporary file")
}

func markdown2xhtml(markdown string) []byte {
	flags := 0 |
		blackfriday.HTML_USE_XHTML
	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS
	renderer := blackfriday.HtmlRenderer(flags, "", "")
	output := blackfriday.Markdown([]byte(markdown), renderer, extensions)
	return []byte("<xhtml>\n<body>\n\n" + string(output) + "\n</body>\n</xhtml>")
}

func markdownData(text string) (MetaData, string) {
	var data MetaData
	r := regexp.MustCompile("(?ms)\\A---.*?(---|\\.\\.\\.)")
	yml := r.FindString(text)
	if yml == "" {
		return data, text
	}
	err := yaml.Unmarshal([]byte(yml), &data)
	printError(err, "Error unmarshalling YAML markdown data")
	return data, text[len(yml):]
}

func imageDir(text, imgDir string) string {
	absDir, err := filepath.Abs(imgDir)
	printError(err, "Error getting absolute file path")
	r := regexp.MustCompile(`!\[(.*?)\]\(([^/].*?/)*([^/]*?)\)`)
	return r.ReplaceAllString(text, "![$1]("+absDir+"/$3)")
}

func absoluteDir(text, filename string) string {
	absDir, err := filepath.Abs(filepath.Dir(filename))
	printError(err, "Error getting absolute file path")
	r := regexp.MustCompile(`!\[(.*?)\]\(([^/].*?)\)`)
	return r.ReplaceAllString(text, "![$1]("+absDir+"/$2)")
}

func generatePdf(xhtmlFile, outFile string, data map[string]string) {
	if data["date"] == "" {
		data["date"] = time.Now().Local().Format("20060102")
	}
	if data["lang"] == "" {
		data["lang"] = "fr"
	}
	params := []string{
		data["date"],
		"htmldoc",
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
	command := exec.Command("faketime", params...)
	env := os.Environ()
	for i, e := range env {
		if strings.HasPrefix(e, "LANG") {
			env[i] = "LANG=" + LOCALE[data["lang"]]
		}
		if strings.HasPrefix(e, "LC_ALL") {
			env[i] = "LC_ALL=" + LOCALE[data["lang"]]
		}
	}
	command.Env = env
	executeCommand(command)
}

func processFile(filename string, printXhtml, printHtml bool, imgDir, outFile string) {
	source, err := ioutil.ReadFile(filename)
	printError(err, "Error reading file")
	data, markdown := markdownData(string(source))
	if imgDir != "" {
		markdown = imageDir(markdown, imgDir)
	} else {
		markdown = absoluteDir(markdown, filename)
	}
	xhtml := markdown2xhtml(markdown)
	if printXhtml {
		fmt.Println(string(xhtml))
		return
	}
	tmpFile, err := ioutil.TempFile("/tmp", "md2pdf-")
	printError(err, "Error creating temporary file")
	defer os.Remove(tmpFile.Name())
	ioutil.WriteFile(tmpFile.Name(), xhtml, 0644)
	processXsl(tmpFile.Name(), data.ToMap())
	if printHtml {
		source, err := ioutil.ReadFile(tmpFile.Name())
		printError(err, "Error reading temporary file")
		fmt.Println(string(source))
		return
	}
	if len(outFile) == 0 {
		outFile = filename[0:len(filename)-len(filepath.Ext(filename))] + ".pdf"
	}
	generatePdf(tmpFile.Name(), outFile, data.ToMap())
}

func main() {
	xhtml := false
	html := false
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
		} else if arg == "-t" || arg == "--html" {
			html = true
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
	processFile(file, xhtml, html, imgDir, outFile)
}
