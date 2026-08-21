package models

import (
	"fmt"
	"os"
	"testing"
)

func TestParseSourceBody(t *testing.T) {

	var index Website
	dat, err := os.ReadFile("index.html")
	if err != nil {
		t.Fatalf("Error reading index.html")
	}

	index.Body = string(dat)

	subpages, err := ParseSourceBody(&index)
	if err != nil {
		t.Fatalf("Error reading index.html")
	}

	if len(subpages) == 0 {
		t.Fatalf("no websites")
	}

	for _, s := range subpages {
		fmt.Printf("%s\n", s.URL)
	}

}
