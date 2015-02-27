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
	HELP = `md2xml [-h] [-x] [-s] [-a] [-i dir] [-o file] file.md
Transform a given Markdown file into XML.
-h        To print this help page.
-x        Print intermediate XHTML output.
-s        Print stylesheet used for transformation.
-a        Output article (instead of blog entry).
-i dir    To indicate image directory.
-o file   The name of the file to output.
file.md   The markdown file to convert.
Note: this program calls pandoc and xsltproc that must have been installed.`
	STYLESHEET_ARTICLE = `<?xml version="1.0" encoding="utf-8"?>

<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                version="1.0">

  <xsl:output method="xml" encoding="UTF-8"/>
  <xsl:param name="id">ID</xsl:param>
  <xsl:param name="date">DATE</xsl:param>
  <xsl:param name="title">TITLE</xsl:param>
  <xsl:param name="author">AUTHOR</xsl:param>
  <xsl:param name="email">EMAIL</xsl:param>
  <xsl:param name="lang">fr</xsl:param>
  <xsl:param name="toc">yes</xsl:param>

  <!-- catch the root element -->
  <xsl:template match="/xhtml">
    <xsl:text disable-output-escaping="yes">
    &lt;!DOCTYPE article PUBLIC "-//CAFEBABE//DTD blog 1.0//EN"
                             "../dtd/article.dtd">
    </xsl:text>
    <article>
      <xsl:attribute name="id"><xsl:value-of select="$id"/></xsl:attribute>
      <xsl:attribute name="date"><xsl:value-of select="$date"/></xsl:attribute>
      <xsl:attribute name="author"><xsl:value-of select="$author"/></xsl:attribute>
      <xsl:attribute name="email"><xsl:value-of select="$email"/></xsl:attribute>
      <xsl:attribute name="lang"><xsl:value-of select="$lang"/></xsl:attribute>
      <xsl:attribute name="toc"><xsl:value-of select="$toc"/></xsl:attribute>
      <title><xsl:value-of select="$title"/></title>
      <text>
       <xsl:apply-templates/>
      </text>
    </article>
  </xsl:template>

  <xsl:template match="h1">
    <sect level="1"><title><xsl:value-of select="."/></title></sect>
  </xsl:template>

  <xsl:template match="h2">
    <sect level="2"><title><xsl:value-of select="."/></title></sect>
  </xsl:template>

  <xsl:template match="h3">
    <sect level="3"><title><xsl:value-of select="."/></title></sect>
  </xsl:template>

  <xsl:template match="h4">
    <sect level="4"><title><xsl:value-of select="."/></title></sect>
  </xsl:template>

  <xsl:template match="h5">
    <sect level="5"><title><xsl:value-of select="."/></title></sect>
  </xsl:template>

  <xsl:template match="h6">
    <sect level="6"><title><xsl:value-of select="."/></title></sect>
  </xsl:template>

  <xsl:template match="p[@class='caption']">
  </xsl:template>

  <xsl:template match="p[count(text())=0 and count(code)=1]">
    <source><xsl:apply-templates select="code"/></source>
  </xsl:template>

  <xsl:template match="p[count(text())=1 and count(img)=1]">
    <xsl:apply-templates select="img"/>
  </xsl:template>

  <xsl:template match="img">
    <figure url="{@src}">
      <xsl:if test="@title">
        <title><xsl:value-of select="@title"/></title>
      </xsl:if>
    </figure>
  </xsl:template>

  <xsl:template match="p">
    <p><xsl:apply-templates/></p>
  </xsl:template>

  <xsl:template match="ul">
    <list><xsl:apply-templates/></list>
  </xsl:template>

  <xsl:template match="ol">
    <enum><xsl:apply-templates/></enum>
  </xsl:template>

  <xsl:template match="li">
    <item><xsl:apply-templates/></item>
  </xsl:template>

  <xsl:template match="table">
    <table><xsl:apply-templates/></table>
  </xsl:template>

  <xsl:template match="tr[count(th)=0]">
    <li><xsl:apply-templates/></li>
  </xsl:template>

  <xsl:template match="tr[count(th) &gt; 0]">
    <th><xsl:apply-templates/></th>
  </xsl:template>

  <xsl:template match="th">
    <co><xsl:apply-templates/></co>
  </xsl:template>

  <xsl:template match="td">
    <co><xsl:apply-templates/></co>
  </xsl:template>

  <xsl:template match="pre">
    <source><xsl:apply-templates/></source>
  </xsl:template>

  <xsl:template match="em">
    <term><xsl:apply-templates/></term>
  </xsl:template>

  <xsl:template match="strong">
    <imp><xsl:apply-templates/></imp>
  </xsl:template>

  <xsl:template match="a">
    <link url="{@href}"><xsl:apply-templates/></link>
  </xsl:template>

</xsl:stylesheet>`
	STYLESHEET_BLOG = `<?xml version="1.0" encoding="utf-8"?>

<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                version="1.0">

  <xsl:output method="xml" encoding="UTF-8"/>
  <xsl:param name="id">ID</xsl:param>
  <xsl:param name="date">DATE</xsl:param>
  <xsl:param name="title">TITLE</xsl:param>

  <!-- catch the root element -->
  <xsl:template match="/xhtml">
    <xsl:text disable-output-escaping="yes">
    &lt;!DOCTYPE blog PUBLIC "-//CAFEBABE//DTD blog 1.0//EN"
                          "../dtd/blog.dtd">
    </xsl:text>
    <blog>
      <xsl:attribute name="id"><xsl:value-of select="$id"/></xsl:attribute>
      <xsl:attribute name="date"><xsl:value-of select="$date"/></xsl:attribute>
      <title><xsl:value-of select="$title"/></title>
      <xsl:apply-templates/>
    </blog>
  </xsl:template>

  <xsl:template match="h1">
    <p><imp><xsl:value-of select="."/></imp></p>
  </xsl:template>

  <xsl:template match="h2">
    <p><imp><xsl:value-of select="."/></imp></p>
  </xsl:template>

  <xsl:template match="h3">
    <p><imp><xsl:value-of select="."/></imp></p>
  </xsl:template>

  <xsl:template match="h4">
    <p><imp><xsl:value-of select="."/></imp></p>
  </xsl:template>

  <xsl:template match="h5">
    <p><imp><xsl:value-of select="."/></imp></p>
  </xsl:template>

  <xsl:template match="h6">
    <p><imp><xsl:value-of select="."/></imp></p>
  </xsl:template>

  <xsl:template match="p[@class='caption']">
  </xsl:template>

  <xsl:template match="p[count(text())=0 and count(code)=1]">
    <source><xsl:apply-templates select="code"/></source>
  </xsl:template>

  <xsl:template match="p[count(text())=1 and count(img)=1]">
    <xsl:apply-templates select="img"/>
  </xsl:template>

  <xsl:template match="img">
    <figure url="{@src}">
      <xsl:if test="@title">
        <title><xsl:value-of select="@title"/></title>
      </xsl:if>
    </figure>
  </xsl:template>

  <xsl:template match="p">
    <p><xsl:apply-templates/></p>
  </xsl:template>

  <xsl:template match="ul">
    <list><xsl:apply-templates/></list>
  </xsl:template>

  <xsl:template match="ol">
    <enum><xsl:apply-templates/></enum>
  </xsl:template>

  <xsl:template match="li">
    <item><xsl:apply-templates/></item>
  </xsl:template>

  <xsl:template match="table">
    <table><xsl:apply-templates/></table>
  </xsl:template>

  <xsl:template match="tr[count(th)=0]">
    <li><xsl:apply-templates/></li>
  </xsl:template>

  <xsl:template match="tr[count(th) &gt; 0]">
    <th><xsl:apply-templates/></th>
  </xsl:template>

  <xsl:template match="th">
    <co><xsl:apply-templates/></co>
  </xsl:template>

  <xsl:template match="td">
    <co><xsl:apply-templates/></co>
  </xsl:template>

  <xsl:template match="pre">
    <source><xsl:apply-templates/></source>
  </xsl:template>

  <xsl:template match="em">
    <term><xsl:apply-templates/></term>
  </xsl:template>

  <xsl:template match="strong">
    <imp><xsl:apply-templates/></imp>
  </xsl:template>

  <xsl:template match="a">
    <link url="{@href}"><xsl:apply-templates/></link>
  </xsl:template>

</xsl:stylesheet>`
	XHTML_HEADER = "<xhtml>\n"
	XHTML_FOOTER = "\n</xhtml>"
)

func processXsl(xmlFile string, data map[string]string, article bool) []byte {
	xslFile, err := ioutil.TempFile("/tmp", "md2xsl-")
	if err != nil {
		panic(err)
	}
	stylesheet := STYLESHEET_ARTICLE
	if !article {
		stylesheet = STYLESHEET_BLOG
	}
	err = ioutil.WriteFile(xslFile.Name(), []byte(stylesheet), 0644)
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
	params = append(params, xmlFile)
	command := exec.Command("xsltproc", params...)
	result, err := command.CombinedOutput()
	if err != nil {
		println(result)
		panic(err)
	}
	return result
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

func processFile(filename string, printXhtml bool, article bool, imgDir, outFile string) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	data, markdown := markdownData(string(source))
	markdown = imageDir(markdown, imgDir)
	xhtml := markdown2xhtml(markdown)
	if printXhtml {
		fmt.Println(string(xhtml))
	}
	xmlFile, err := ioutil.TempFile("/tmp", "md2xml-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(xmlFile.Name())
	ioutil.WriteFile(xmlFile.Name(), xhtml, 0644)
	result := processXsl(xmlFile.Name(), data, article)
	if len(outFile) > 0 {
		ioutil.WriteFile(outFile, result, 0644)
	} else {
		fmt.Println(string(result))
	}
}

func main() {
	xhtml := false
	article := false
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
		} else if arg == "-a" || arg == "--article" {
			article = true
		} else if arg == "-s" || arg == "--stylesheet" {
			if article {
				fmt.Println(STYLESHEET_ARTICLE)
			} else {
				fmt.Println(STYLESHEET_BLOG)
			}
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
	processFile(file, xhtml, article, imgDir, outFile)
}
