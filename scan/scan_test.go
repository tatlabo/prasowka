package scan

import (
	"html/template"
	"prasowka/internal/models"
	"testing"
	"time"
)

var path = "https://www.rmf24.pl/"

func TestScrapPage(t *testing.T) {

	start := time.Now()
	defer func() {
		t.Logf("Scan completed in %s\n", time.Since(start))
	}()

	w := models.Website{}

	w.URL = template.URL(path)

	if err := w.ProcessWebsite(); err != nil {
		t.Errorf("Error processing website: %v", err)
	}

	err := ScrapPage(w)
	if err != nil {
		t.Errorf("%v", ErrScaningPage)
	}

}
