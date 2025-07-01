package main

import (
	"fmt"
	"job-pulse/internal/config"
)

func main() {
	fmt.Println("Hello, job-pulse!")
	cfg := config.MustLoad()

	fmt.Println(cfg)
}