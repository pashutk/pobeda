package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Poehali!")

	to := "LCA"
	if len(os.Args) > 1 {
		to = os.Args[1]
	}

	flights := getFlightsForRegion(to)

	for _, flight := range flights {
		fmt.Println(flight)
	}

	fmt.Println("Pobeda!")
}
