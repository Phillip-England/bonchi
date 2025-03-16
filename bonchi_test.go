package bonchi

import (
	"testing"
)

func TestBonchi(t *testing.T) {
	_, err := BundleCss("./css", "./output.css")
	if err != nil {
		panic(err)
	}
	_, err = BundleJs("./js", "./output.js")
	if err != nil {
		panic(err)
	}
}
