# bonchi
A minimal web bundler with css preprocessing and js support.

## Installation
```bash
go get github.com/Phillip-England/bonchi
```

## Hello, World!
```go
package main

import (
  "fmt"

  "github.com/Phillip-England/bonchi"
)

func main() {
  css, err := bonchi.BundleCss("./css", "./static/output.css")
  if err != nil {
    panic(err)
  }
  js, err := bonchi.BundleJs("./js", "./static/output.js")
  if err != nil {
    panic(err)
  }
  fmt.Println(css, js)
}
```

## Target Dir
Bonchi is based off of a target directory. All the files in the directory will be bundled and processed for mixin support.

## File Names and Ordering
File names can be used to order the way files are organized in the output. For example, `0.reset.css` will be displayed first in the output, followed by `1.header.css` and so on. The same goes for `.js` files as well.

## Mixin
Any class can be used within another class using `bonchi-mix:".className1 .className2";`

```css
.blue {
  background-color:skyblue;
}

.border {
  border:solid black 1px;
}

button {
  bonchi-mix:".blue .border";
}
```