package main

func main() {
	v, err := Get("file", "Inspect.mp4")
	if err != nil {
		panic(err)
	}
}
