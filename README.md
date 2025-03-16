# bonchi
A minimal css preprocessor with importing and mixin support.

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
  css, err := bonchi.Bundle("./css", "./static/output.css")
  if err != nil {
    panic(err)
  }
  fmt.Println(css)
}
```

## Target Dir
Bonchi is based off of a target directory. All the files in the directory will be bundled and processed for mixin support.

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