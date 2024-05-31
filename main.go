package main

import "fmt"

func main() {
	kv, err := NewKVStore[string]("kvstore")
	if err != nil {
		fmt.Printf("Error initializing store: %v\n", err)
		return
	}

	// kv.Set("foo", "bar")
	// kv.Set("hello", "world")

	if value, errGet := kv.Get("foo"); errGet == nil {
		fmt.Printf("foo = %s\n", value)
	} else {
		fmt.Println("foo not found")
	}

	if value, errGet := kv.Get("hello"); errGet == nil {
		fmt.Printf("hello = %s\n", value)
	} else {
		fmt.Println("hello not found")
	}
}
