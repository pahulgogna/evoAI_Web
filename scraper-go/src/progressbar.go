package main

import "fmt"

func displayProgress(done int, total int, finished bool) {
	barLen := ((done * 100)/ total)
	fmt.Printf("\r[")

	for i := range 100 {
		if i <= barLen {
			fmt.Printf("=")
		} else {
			fmt.Printf("_")
		}
	}
	fmt.Printf("] %d%%", barLen)

	if finished {
		fmt.Println()
	}
}