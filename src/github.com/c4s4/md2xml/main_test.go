package main

import (
	"testing"
)

func TestTrue(t *testing.T) {
	if true != true {
		t.Error("Impossible failure!")
	}
}

