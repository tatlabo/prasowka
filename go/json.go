package main

import (
	"encoding/json"
	"fmt"
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
}
