MDTOOLS
=======

MdTools is a collection of personal Markdown tools :

- md2pdf : transforms a markdown file into PDF.
- md2xml : transforms a markdown file into XML for my site DTD.

They both call *pandoc* and *xsltproc* that must have been installed. *md2pdf*
also calls *htmldoc* which is a tool to transforms HTML files to PDF and 
*faketime* to fake date while generating PDF. To install these dependencies on
a Debian like Linux distribution, you could type following commands :

    sudo apt-get install pandoc
	sudo apt-get install xsltproc
	sudo apt-get install htmldoc
    sudo apt-get install faketime

md2pdf
------

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

### Note

I known that *pandoc* is supposed to produce PDF output with command such as:

     $ pandoc -f markdown -t latex -o example.pdf example.md

But I am not satisfied with resulting output.

md2xml
------

To get help about this tool, type :

    $ md2xml -h
	md2xml [-h] [-x] [-s] [-a] [-i dir] [-o file] file.md
    Transform a given Markdown file into XML.
    -h        To print this help page.
    -x        Print intermediate XHTML output.
    -s        Print stylesheet used for transformation.
    -a        Output article (instead of blog entry).
    -p        Add link to PDF version.
    -i dir    To indicate image directory.
    -o file   The name of the file to output.
    file.md   The markdown file to convert.
    Note: this program calls pandoc and xsltproc that must have been installed.

This tool calls *pandoc* to transform the markdown file into an XHTML one. This
is the file printed with the `-x` option. This file is transformed, calling 
*xsltproc* and the stylesheet printed with the ̀`-s` option, into the resulting
XML file. The `-a` option prints an article instead of an blog entry (using a
different stylesheet).

The option `-i dir` tells in which directory are located images (relative to
current directory).

The `-p` option asks for a link to the PDF file into resulting HTML file.

This will print resulting XML document on the console. To put result into a
given file, use the `-o` option.

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

Releases
--------

- **1.1.1** (2015-03-22): Added faketime dependency in documentation.
- **1.1.0** (2015-03-22): Fixed date in PDF document footer.
- **1.0.0** (2015-03-17): First release.
