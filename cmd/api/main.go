package main

import (
	"ds-easy/internal/server"
	"fmt"
)

func main() {

	server, servee := server.NewServer()

	servee.TestDB()
	err := server.ListenAndServe()

	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
