package main

import "os"

func notMain() {
	os.Exit(1)
}

func main() {
	os.Exit(1) // want "os.Exit calls are prohibited in main()"
}
