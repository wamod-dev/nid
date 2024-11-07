# nid

[![License][license.icon]][license.page]
[![CI][ci.icon]][ci.page]
[![Coverage][coverage.icon]][coverage.page]
[![Report][report.icon]][report.page]
[![Documentation][docs.icon]][docs.page]
[![Release][release.icon]][release.page]

`nid` package provides human-readable named unique identifiers.

## Overview

The package offers the following functionalities:

- **Named prefix**: All identifiers have name prefix that makes them more readable.
- **Universally unique**: Random part of identifier base is 8 bytes length that is comparable to commonly used UUIDv4 & ULID.
- **Sortable by time**: All identifiers have 8 bytes length time prefix. They are sortable both in string & binary format and have sequential order by time.
- **Text, JSON & SQL Support**: Identifiers implement multiple encoding reducing code needed to add manual conversion:
    -  `encoding.TextMarshaler`, `encoding.TextUnmarshaler` for text encoding.
    -  `json.Marshaler`, `json.Unmarshaler` for JSON encoding.
    -  `sql.Scanner`, `driver.Valuer` for storing in SQL database.

## Installation

To use this package in your Go project, you can import it using:

```go
import "go.wamod.dev/nid"
```

## Usage

To start creating named identifiers, create new `Naming` first:

```go
var BookIDN = nid.MustNaming("book") // snake_case 
```

Create your resource type:

```go
type Book struct {
    ID     nid.NID `json:"id"`
    Title  string  `json:"title"`
    Author string  `json:"author"`
}
```

Then whenever you need to generate new identifier for your resource use identifier `Naming`:

```go
book := Book{
    ID:     BookIDN.New(),
    Title:  "The Hobbit",
    Author: "J.R.R. Tolkien",
}
```

### Helpers

#### Parsing strings

To parse named identifier from string use `Parse` function:

```go
bookID, err := nid.Parse("book_000034o5m20uo63o22umrn7kcs")
```

To parse base identifier from string using `ParseBase` function:

```go
base, err := nid.ParseBase("000034o5m20uo63o22umrn7kcs")
if err != nil {
    // handle error
}

// Apply resource naming to identifier
bookID := BookIDN.Apply(base)
```

#### Converting to string

When you need to convert it to string format you can use `String()` method:

```go
fmt.Print(bookID.String()) 
// Output:
// book_000034o5m20uo63o22umrn7kcs
```

Similarly, you can also get identifier base string value:

```go
fmt.Print(bookID.Base().String())
// Output:
// 000034o5m20uo63o22umrn7kcs
```

#### Compare

To check if two named identifiers are the same you can use equal operator:

```go
var (
    a nid.NID
    b nid.NID
)

if a == b {
    fmt.Print("a == b")
}
```

You can also use `Compare` helpers that return `0`, `-1`, `1`:

```go
switch nid.Compare(a, b) {
    case 0:
        fmt.Print("a == b")
    case 1:
        fmt.Print("a > b")
    case -1:
        fmt.Print("a < b")
}
```

Similarly, you can also compare identifier base:
```go
switch nid.CompareBase(a.Base(), b.Base()) {
    case 0:
        fmt.Print("a == b")
    case 1:
        fmt.Print("a > b")
    case -1:
        fmt.Print("a < b")
}
```

#### Sort

When you need to sort multiple identifiers you can use `Sort` helper:

```go
var ids []nid.NID

nid.Sort(ids)
```

Similarly, you can also sort identifier base:
```go
var baseIDs []nid.Base

nid.SortBase(baseIDs)
```

## Contributing

Thank you for your interest in contributing to the `nid` Go library! We welcome and appreciate any contributions, whether they be bug reports, feature requests, or code changes.

If you've found a bug, please [create an issue][issue.page] describing the problem, including any relevant error messages and a minimal reproduction of the issue.

## License

`nid` is licensed under the [MIT License][license.page].

[issue.page]:    https://github.com/wamod-dev/nid/issues/new/choose
[license.icon]:  https://img.shields.io/badge/license-MIT-green.svg
[license.page]:  https://github.com/wamod-dev/nid/blob/main/LICENSE
[ci.icon]:       https://github.com/wamod-dev/nid/actions/workflows/go.yml/badge.svg
[ci.page]:       https://github.com/wamod-dev/nid/actions/workflows/go.yml
[coverage.icon]: https://codecov.io/gh/wamod-dev/nid/graph/badge.svg?token=YO8HCMUVK9
[coverage.page]: https://codecov.io/gh/wamod-dev/nid
[report.icon]:   https://goreportcard.com/badge/go.wamod.dev/nid
[report.page]:   https://goreportcard.com/report/go.wamod.dev/nid
[docs.icon]:     https://godoc.org/go.wamod.dev/nid?status.svg
[docs.page]:     http://godoc.org/go.wamod.dev/nid
[release.icon]:  https://img.shields.io/github/release/wamod-dev/nid.svg
[release.page]:  https://github.com/wamod-dev/nid/releases/latest