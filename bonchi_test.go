package bonchi

import (
	"testing"
)

func TestBonchi(t *testing.T) {
	_, err := Bundle("./css/input.css", "./css/output.css")
	if err != nil {
		panic(err)
	}
}
