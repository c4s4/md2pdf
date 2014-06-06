package main

import (
    "os"
    "github.com/russross/blackfriday"
    "os/exec"
    "io/ioutil"
    "fmt"
)

const stylesheet = `<?xml version="1.0" encoding="utf-8"?>
<!--
Stylesheet to transform an XHTML document to XML one.
-->

<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                version="1.0">

  <xsl:output method="xml" encoding="UTF-8"/>
  <xsl:param name="id">ID</xsl:param>
  <xsl:param name="date">DATE</xsl:param>
  <xsl:param name="title">TITLE</xsl:param>

  <!-- catch the root element -->
  <xsl:template match="/xhtml">
    <xsl:text disable-output-escaping="yes">
    &lt;!DOCTYPE weblog PUBLIC "-//CAFEBABE//DTD weblog 1.0//EN"
                               "../dtd/weblog.dtd">
    </xsl:text>
    <weblog>
      <xsl:attribute name="id"><xsl:value-of select="$id"/></xsl:attribute>
      <xsl:attribute name="date"><xsl:value-of select="$date"/></xsl:attribute>
      <title><xsl:value-of select="$title"/></title>
      <xsl:apply-templates/>
    </weblog>
  </xsl:template>

  <xsl:template match="h1">
    <p><imp><xsl:apply-templates/></imp></p>
  </xsl:template>

  <xsl:template match="h2">
    <p><imp><xsl:apply-templates/></imp></p>
  </xsl:template>

  <xsl:template match="h3">
    <p><imp><xsl:apply-templates/></imp></p>
  </xsl:template>

  <xsl:template match="h4">
    <p><imp><xsl:apply-templates/></imp></p>
  </xsl:template>

  <xsl:template match="h5">
    <p><imp><xsl:apply-templates/></imp></p>
  </xsl:template>

  <xsl:template match="h6">
    <p><imp><xsl:apply-templates/></imp></p>
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

  <xsl:template match="th">
    <th><xsl:apply-templates/></th>
  </xsl:template>

  <xsl:template match="tr">
    <li><xsl:apply-templates/></li>
  </xsl:template>

  <xsl:template match="td">
    <co><xsl:apply-templates/></co>
  </xsl:template>

  <xsl:template match="pre">
    <source><xsl:apply-templates/></source>
  </xsl:template>

  <xsl:template match="img">
    <figure url="{@src}" width="{@width}" height="{@height}"/>
  </xsl:template>

  <xsl:template match="em">
    <term><xsl:apply-templates/></term>
  </xsl:template>

  <xsl:template match="strong">
    <imp><xsl:apply-templates/></imp>
  </xsl:template>

  <xsl:template match="code">
    <code><xsl:apply-templates/></code>
  </xsl:template>

</xsl:stylesheet>`

func processXsl(xmlFile string) ([]byte) {
    xslFile, err := ioutil.TempFile("/tmp", "md2xsl-")
    if err != nil {
        panic(err)
    }
    err = ioutil.WriteFile(xslFile.Name(), []byte(stylesheet), 0755)
    if err != nil {
        panic(err)
    }
    defer os.Remove(xslFile.Name())
    command := exec.Command("xsltproc", xslFile.Name(), xmlFile)
    result, err := command.CombinedOutput()
    if err != nil {
        panic(err)
    }
    return result
}

func markdown2xhtml(filename string) ([]byte) {
    markdown, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    xhtml := "<xhtml>\n" + string(blackfriday.MarkdownCommon([]byte(markdown))) + "\n</xhtml>"
    return []byte(xhtml)
}

func processFile(filename string) string {
    xhtml := markdown2xhtml(filename)
    xmlFile, err := ioutil.TempFile("/tmp", "md2xml-")
    if err != nil {
        panic(err)
    }
    defer os.Remove(xmlFile.Name())
    ioutil.WriteFile(xmlFile.Name(), xhtml, 0755)
    result := processXsl(xmlFile.Name())
    return string(result)
}

func main() {
    for _, filename := range os.Args[1:len(os.Args)] {
        fmt.Println(processFile(filename))
    }
}

