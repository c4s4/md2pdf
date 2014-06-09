MD2XML
======

To transform Markdown document to XML one with Sweetohm DTD.

Usage
-----

To transform document *test.md* to XML, type following command:

```
md2xml test.md
```

This will output resulting XML document to the console.

Markdown syntax
---------------

See file *test/test.md* for supported syntax elements. This is syntax described on [markdown wiki page](http://en.wikipedia.org/wiki/Markdown), plus images with following syntax:

```
![Image alternative](image_file.jpg "Image title")
```

Bugs
----

### Successive lists

If an ordered list follows an unordered one, it result in a single unordered list:

```
- First unordered.
- Second unordered.
- Third unordered.

1. First ordered.
2. Second ordered.
3. Third ordered.
```

If there is a paragraph between, it works:

```
- First unordered.
- Second unordered.
- Third unordered.

Test.

1. First ordered.
2. Second ordered.
3. Third ordered.
```
