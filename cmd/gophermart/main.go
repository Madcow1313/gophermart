package main

import "github.com/Madcow1313/gophermart/internal/server"

func main() {
	s := server.New("localhost:8080")
	s.DatabaseDSN = "host=localhost user=postgres password=postgres"
	s.Serve()
}
