package main

import "fmt"

//go:generate go run generate.go
//go:generate go fmt models.go

func main() {
	fmt.Println("you are on main")
}
