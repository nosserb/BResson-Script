package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 6 {
		return
	}

	for i := 1; i <= 5; i++ {
		arg := os.Args[i]
		if arg == "01" || arg == "galaxy 01" {
			fmt.Println("Alert!!!")
			return
		}
	}
}
