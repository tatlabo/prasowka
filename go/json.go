package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
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
}
