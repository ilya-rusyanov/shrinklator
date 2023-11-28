package main

import "os"

func main() {
	os.Exit(1) // want "os.Exit calls in main() are prohibited"
}
