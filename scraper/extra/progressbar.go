package extra

import "fmt"

func DisplayProgress(done int, total int, finished bool) {
	
	if total == 0 {
        return
    }

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