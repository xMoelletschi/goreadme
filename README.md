# goreadme

Package goreadme generates readme markdown file from go doc.

The package can be used as a command line tool, described below:

# Use as a command line tool

```go
$ GO111MODULE=on go get github.com/xMoelletschi/goreadme/cmd/goreadme
$ goreadme -h
```

# Why Should You Use It

Both Go doc and readme files are important. Go doc to be used by your user's library, and README
file to welcome users to use your library. They share common content, which is usually duplicated
from the doc to the readme or vice versa once the library is ready. The problem is that keeping
documentation updated is important, and hard enough - keeping both updated is twice as hard.

# Go Doc Instructions

The formatting of the README.md is done by the go doc parser. This makes the result README.md a
bit more limited. Currently, `goreadme` supports the formatting as explained in
[godoc page](https://blog.golang.org/godoc-documenting-go-code), or
[here](https://pkg.go.dev/github.com/fluhus/godoc-tricks). Meaning:

* A header is a single line that is separated from a paragraph above.

* Code block is recognized by indentation as Go code.

```go
func main() {
  ...
}
```

* Inline code is marked with `backticks`.

* URLs will just automatically be converted to links: [https://github.com/xMoelletschi/goreadme](https://github.com/xMoelletschi/goreadme)

Additionally, the syntax was extended to include some more markdown features while keeping the Go
doc readable:

* Bulleted and numbered lists are possible when each bullet item is followed by an empty line.

* Diff blocks are automatically detected when each line in a code block starts with a `' '`,
`'-'` or `'+'`:

```diff
-removed line starts with '-'
 remained line starts with ' '
+added line starts with '+'
```

* A repository file can be linked when providing a path that start with `[./](./)`: [./goreadme.go](./goreadme.go).

* A link can have a link text by prefixing it with parenthesised text:
[goreadme page](https://github.com/xMoelletschi/goreadme).

* A link to repository file and can have a link text: [goreadme main file](./goreamde.go).

* An image can be added by prefixing a link to an image with `(image/<image title>)`:

![title of image](https://github.githubassets.com/images/icons/emoji/unicode/1f44c.png)

# Testing

The goreadme tests the test cases in the [./testdata](./testdata) directory. It generates readme files for
all the packages in that directory and asserts that the result readme matches the existing one.
When modifying goreadme behavior, there is no need to manually change these readme files. It is
possible to run `WRITE_READMES=1 go test ./...` which regenerates them and check the changes
match the expected (optionally using `git diff`).

---
Readme created from Go doc with [goreadme](https://github.com/xMoelletschi/goreadme)
