package main

import (
	"ds-easy/internal/server"
	"fmt"
)

func main() {
	fmt.Printf("DSEASY")
	server := server.NewServer()

	err := server.ListenAndServe()

	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
