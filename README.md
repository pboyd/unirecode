# Unicode Character Encoding

This project contains implementations of various Unicode character encodings in
Go and a command line tool to convert between them.

This was written as an exercise not as real code. The goal was to write simple
implementations of common Unicode character encodings. So if you like fiddling
around with character encodings, have a look. But otherwise, you should look
elsewhere:

* For a Go text encoding library, the standard library's [unicode/utf8](https://golang.org/pkg/unicode/utf8/) and [unicode/utf16](https://golang.org/pkg/unicode/utf16/) packages solve a lot of common problems.
* If you need more than that, try [golang.org/x/text](https://godoc.org/golang.org/x/text).
* For a command line program to translate between character encodings, use `iconv`. You probably have installed already.
