package main

import "fmt"

func main() {
	err := server()

	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
