package rssnews

import (
	"strings"
	"testing"
)

func TestUnmarshalFeed(t *testing.T) {
	feedURL := "https://www.rmf24.pl/feed"
	data, err := ReadRRSS(feedURL)
	if err != nil {
		t.Fatal(err)
	}

	feed, err := UnmarshalFeed(data)
	if err != nil {
		t.Fatal(err)
	}

	if feed.Version != "2.0" {
		t.Fatalf("version = %q, want %q", feed.Version, "2.0")
	}
	if feed.Channel.Title != "Najnowsze wiadomości - RMF24.pl" {
		t.Fatalf("channel title = %q", feed.Channel.Title)
	}
	if len(feed.Channel.Items) == 0 {
		t.Fatal("expected feed items")
	}

	item := feed.Channel.Items[0]
	if item.GUID != "1015925" {
		t.Fatalf("first item GUID = %q, want %q", item.GUID, "1015800")
	}
	if item.Enclosure.Type != "image/jpeg" {
		t.Fatalf("enclosure type = %q, want %q", item.Enclosure.Type, "image/jpeg")
	}
	if !strings.Contains(item.Description, "było celowe podpalenie? Nowe informacje ws.") {
		t.Fatalf("description = %q, want feed content", item.Description)
	}
}
