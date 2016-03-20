MD2PDF
======

md2pdf is a tool to convert [Markdown](https://en.wikipedia.org/wiki/Markdown)
documents to PDF, without using LaTeX.

Installation
------------

Download the archive on 
[the releases page](https://github.com/c4s4/md2pdf/releases). Unzip it and
put the binary for your platform somewhere in your *PATH* (in directory
*/usr/local/bin* for instance).

As this tool calls *pandoc*, *htmldoc*, *xsltproc* and *faketime*, you must
install them with *md2pdf*. To install these dependencies on a Debian like
Linux distribution, you can type following commands :

    sudo apt-get install pandoc
	sudo apt-get install xsltproc
	sudo apt-get install htmldoc
    sudo apt-get install faketime

Usage
-----

To get help about this tool, type :

    $ md2pdf -h
    md2pdf [-h] [-x] [-s] [-t] [-i dir] [-o file] file.md
    Transform a given Markdown file into PDF.
    -h        To print this help page.
    -x        Print intermediate XHTML output.
    -s        Print stylesheet used for transformation.
    -t        Print html output.
    -i dir    To indicate image directory.
    -o file   The name of the file to output.
    file.md   The markdown file to convert.
    Note:
    This program calls pandoc, xsltproc and htmldoc that must have been installed.

This tool calls *pandoc* to transform the markdown file into an XHTML one. This
is the file printed with the `-x` option. This file is transformed, calling 
*xsltproc* and the stylesheet printed with the ̀`-s` option, into an intermediate
decorated XHTML file printed with the `-t` option. This file is transformed into
resulting PDF calling *xsltproc*.

The option `-i dir` tells in which directory are located images (relative to
current directory).

This will print resulting PDF document in a file with the same path than the
origin markdown document with the *.pdf* extension. To write PDF in another file
use the `-o file` option.

Markdown syntax
---------------

See file *test/example.md* for an example of supported syntax elements. This is
syntax described on [markdown wiki page](http://en.wikipedia.org/wiki/Markdown),
plus images with following syntax:

    ![Image Title](image_file.jpg)

Furthermore, these tools parse special information headers at the beginning of
the markdown files, such as :

    % id:       1
    % date:     2014-06-09
    % title:    Document title
    % author:   Michel Casabianca
    % email:    michel.casabianca@gmail.com
	% lang:     en
	% toc:      true
    % keywords: markdown test

These headers are used by tools:

- **id**: this is the unique ID of the document. This is for my site processor
  internal usage.
- **date**: this ISO formatted date is print in documents as production date.
- **title**: this is the title of the document.
- **auhtor**: this is the author of the document.
- **email**: this is the email of the author of the document.
- **lang**: this is the language of the document, in ISO format (defaults to
  *en*).
- **toc**: tells if we want a table of content in resulting document (defaults
  to *false*).
- **keywords**: this is for the internal usage of the site generator.

Note that these headers are not mandatory and will have default values if not
set.

Bugs
----

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

Todo
----

- Enable TOC for PDF.

