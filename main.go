package main

import (
	"fmt"
	"github.com/SyntaxSamurai/Bootdev/BlogAggregator/internal/config"
)

func main() {
	res, err := config.Read()
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}
	fmt.Println("values received: ======== ",res)
}
