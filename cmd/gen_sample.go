package main

import (
	"fmt"
)

func main() {
	v, err := GetVideo("file", "newrising.mp4")
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
}
