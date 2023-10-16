package utils

import (
	"fmt"
	"strconv"
)

// print a single process bar
func printBar(done int, load int, pref string) {
	if load == 0 || done > load {
		return
	}

	suff := " (" + strconv.Itoa(done) + "/" + strconv.Itoa(load) + ")"
	if done == load {
		suff += " Success!"
	}

	ratio := done * 50 / load
	fmt.Printf(pref + "[")
	for i := 0; i < 50; i++ {
		if i <= ratio {
			fmt.Printf("=")
		} else {
			fmt.Printf("-")
		}
	}
	fmt.Println("]" + suff)
}

// print multi or a single process bar
func ProgressBar(n int, dones []int, loads []int, prefs []string) {
	if len(dones) == n && len(loads) == n && len(prefs) == n {
		fmt.Printf("\033[?25l\033[%dA\033[K", n)
		for i := 0; i < n; i++ {
			printBar(dones[i], loads[i], prefs[i])
		}
	}
}
