# bonchi
A minimal css preprocessor with importing and mixin support.

## Installation
```bash
go get github.com/Phillip-England/bonchi
```

## Hello, World!
To get started, you'll need three files `./input.css`, `./global.css`, and `main.go`.

In `./input.css`, we use `@bonchi` to import css and `bonchi-mix` to mixin the css from other classes.

**BIG GOTCHA HERE**, the import you use with `@bonchi` is relative to the root of your **GO PROGRAM**, not the **INPUT FILE** itself.
```css
@bonchi ./global.css;

.some-element {
  bonchi-mix:".hidden .bg-blue";
}
```

In `./global.css`, we define some css to be imported.
```css
html, body {
  min-height: 100vh;
}

* {
  margin:0;
  padding:0;
  box-sizing: border-box;
}

.hidden {
  visibility: hidden;
}

.bg-blue {
  background-color:blue;
}
```

Once you get your css lined up, in `main.go` run:
```go
package main

import (
  "fmt"

  "github.com/Phillip-England/bonchi"
)

func main() {
  css, err := bonchi.Bundle("./input.css", "./output.css")
  if err != nil {
    panic(err)
  }
  fmt.Println(css)
}
```