package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// https://youtu.be/T2fqLam1iuk?t=860

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type response struct {
	Item   string `json:"item"`
	Album  string
	Title  string
	Artist string
}

type respWrapper struct {
	response
}

var j1 = `{
	"item": "album",
	"album": {"title": "The Dark Side of the Moon"}
}`

var j2 = `{
	"item": "song",
	"song": {"title": "The great gig in the sky", "artist": "Pink Floyd"}
}`

func (p Person) String() string {
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}

func main() {
	// helper()
	var rw1, rw2 respWrapper

	err := rw1.UnmarshalJSON([]byte(j1))
	if err != nil {
		log.Fatalf("Error unmarshaling j1: %v", err)
	}

	err = rw2.UnmarshalJSON([]byte(j2))
	if err != nil {
		log.Fatalf("Error unmarshaling j2: %v", err)
	}

}

func (rw *respWrapper) UnmarshalJSON(b []byte) (err error) {
	var raw map[string]any

	err = json.Unmarshal(b, &rw.response)
	err = json.Unmarshal(b, &raw)

	log.Printf("Unmarshaled response: %+v", rw.response)

	return nil
}

func helper() {
	alice := Person{Name: "Alice", Age: 30}

	out, err := json.Marshal(&alice)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Printf("%v\n", string(out))

	d := "2026-01-05T09:45:36"

	t, err := time.Parse("2006-01-02T15:04:05", d)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
	fmt.Println("Parsed time:", t)

	i := any("example string")

	log.Printf("Value: %v, Type: %T", i, i)

	log.Printf("Person: %v, Type: %T", alice, alice)

}
