# MD2PDF

<!--
[![Build Status](https://travis-ci.org/c4s4/md2pdf.svg?branch=master)](https://travis-ci.org/c4s4/md2pdf)
-->
[![Code Quality](https://goreportcard.com/badge/github.com/c4s4/md2pdf)](https://goreportcard.com/report/github.com/c4s4/md2pdf)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
<!--
[![Coverage Report](https://coveralls.io/repos/github/c4s4/md2pdf/badge.svg?branch=master)](https://coveralls.io/github/c4s4/md2pdf?branch=master)
-->

- Project :   <https://github.com/c4s4/md2pdf>.
- Downloads : <https://github.com/c4s4/md2pdf/releases>.

md2pdf is a tool to convert [Markdown](https://en.wikipedia.org/wiki/Markdown)
documents to PDF, without using LaTeX.

## Installation

### Prerequisites

As this tool calls *htmldoc*, *xsltproc* and *faketime*, you must install them
with *md2pdf*. To install these dependencies on a Debian like Linux
distribution, you can type following commands :

	sudo apt-get install xsltproc
	sudo apt-get install htmldoc
    sudo apt-get install faketime

### Unix users (Linux, BSDs and MacOSX)

Unix users may download and install latest *md2pdf* release with command:

```bash
sh -c "$(curl https://sweetohm.net/dist/md2pdf/install)"
```

If *curl* is not installed on you system, you might run:

```bash
sh -c "$(wget -O - https://sweetohm.net/dist/md2pdf/install)"
```

**Note:** Some directories are protected, even as *root*, on **MacOSX** (since *El Capitan* release), thus you can't install *md2pdf* in */usr/bin* for instance.

### Binary package

Download the archive on
[the releases page](https://github.com/c4s4/md2pdf/releases). Unzip it and
put the binary for your platform somewhere in your *PATH* (in directory
*/usr/local/bin* for instance).

## Usage

To get help about this tool, type :

    $ md2pdf -h
    md2pdf [-h] [-v] [-x] [-s] [-t] [-i dir] [-o file] file.md
    Transform a given Markdown file to PDF:
    -h        To print this help page
    -v        Print version and exit
    -x        Print intermediate XHTML output
    -s        Print stylesheet used for transformation
    -t        Print html output
    -i dir    To indicate image directory
    -o file   The name of the file to output
    file.md   The markdown file to convert
    Note:
    This program calls xsltproc, htmldoc and faketime that must have been installed

This tool transforms Markdown input to XHTML using blackfriday library. This
is the file printed with the `-x` option. This file is transformed, calling
*xsltproc* and the stylesheet printed with the ̀`-s` option, into an intermediate
decorated XHTML file printed with the `-t` option. This file is transformed into
resulting PDF calling *xsltproc*.

The option `-i dir` tells in which directory are located images (relative to
current directory).

This will print resulting PDF document in a file with the same path than the
origin markdown document with the *.pdf* extension. To write PDF in another file
use the `-o file` option.

## Markdown syntax

See file *test/example.md* for an example of supported syntax elements. This is
syntax described on [markdown wiki page](http://en.wikipedia.org/wiki/Markdown),
plus images with following syntax:

    ![Image Title](image_file.jpg)

Furthermore, this tool parses YAML header at the beginning of the markdown
document, as used by *pandoc* tool, such as :

    ---
    title:    Document title
    author:   Michel Casabianca
    date:     2014-06-09
    email:    michel.casabianca@gmail.com
    id:       1
	lang:     en
	toc:      true
    logo:     logo.png
    ---

These headers are used by the tools to print information at the beginning of
the document and in page footer:

- **title**: this is the title of the document.
- **auhtor**: this is the author of the document.
- **date**: this ISO formatted date is print in documents as production date.
- **email**: this is the email of the author of the document.
- **id**: this is the unique ID of the document. This is for my site processor
  internal usage.
- **lang**: this is the language of the document, in ISO format (defaults to
  *en*).
- **toc**: tells if we want a table of content in resulting document (defaults
  to *false*).
- **logo**: logo image inserted in header on the right.

Note that these headers are not mandatory.

## Bugs

### Successive lists

If an ordered list follows an unordered one, it result in a single unordered
list:

    - First unordered.
    - Second unordered.
    - Third unordered.

    1. First ordered.
    2. Second ordered.
    3. Third ordered.

If there is a paragraph between, it works:

    - First unordered.
    - Second unordered.
    - Third unordered.

    Test.

    1. First ordered.
    2. Second ordered.
    3. Third ordered.

## Todo

- Generate TOC in resulting PDF file.
