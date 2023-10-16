package main

import (
	"fmt"
	"os"
	c "pagerank/calcu"
	p "pagerank/parse"
	u "pagerank/utils"
	"time"
)

var N int = len(p.Id2PageMap)
var MT c.Matrix
var PR map[int]float64

func PageRank(pr *map[int]float64, mt *c.Matrix, n int) {
	// record the pagerank's time cost
	fmt.Printf("Start the PageRank Caculate...\n")
	time1 := time.Now()
	defer func() {
		time2 := time.Now()
		cost := time2.Sub(time1).Seconds()
		fmt.Printf("Program PageRank exited. Cost %.2fs\n\n", cost)
	}()

	// intialize the parameters
	var maxC int = 100
	var alpha float64 = 0.85
	var e float64 = (1 - alpha) / float64(N)
	var t float64 = 1e-3

	// save the history value to check if converge
	cp := make(map[int]float64)

	// enter the iterlations
	for i := 0; i < maxC; i++ {
		for k, v := range *pr {
			cp[k] = v
		}

		// PR = alpha * MT * PR + e
		*pr = mt.Mult(pr, true)
		c.VecMult(pr, alpha)
		c.VecAdd(pr, e)

		// show the difference with last time
		diff := c.VecComp(pr, &cp)
		fmt.Printf("\033[?25l\033[%dA\033[KDifference: %.3f%%\n", 1, diff*100)

		// check if converge
		if diff < t {
			fmt.Printf("\nBelow the diffrence threshold. Now existing...\n")
			return
		}
	}
	fmt.Printf("\nReach the max iterations allowed. Now existing...\n")
}

func init() {
	// check if adjax.json exists, parse wiki if not exists
	_, err := os.Stat("adjax.json")
	if err != nil {
		fmt.Println("Can't find adjax file. Now rebuild it...")
		p.ParseWiki()
	}
	fmt.Printf("Find adjax file successfully. Now load it...")

	// initialize the transition matrix
	MT.Load("adjax.json")
	fmt.Printf("OK!\nNow normalize it...")
	MT.Normal(true)
	fmt.Printf("OK!\n\n")

	// initialize the PR vector
	fmt.Printf("Initialize the PR vector...")
	PR = make(map[int]float64, N)
	for id := range p.Id2PageMap {
		PR[id] = 1.0 / float64(N)
	}
	fmt.Printf("OK!\n\n")
}

func main() {
	fmt.Print(len(MT.Col))
	// do the PageRank
	PageRank(&PR, &MT, N)

	// sort the PR vector and get the rank
	fmt.Printf("Now sort the PR vector...")
	rank := u.SortMapByValue(&PR, false)
	fmt.Printf("OK!\n\n")

	// save the PR vector to the rank.txt
	fmt.Printf("save the PR vector to the rank.txt...")
	fil, _ := os.OpenFile("rank.txt", os.O_WRONLY|os.O_CREATE, 0644)
	defer fil.Close()
	for _, r := range rank {
		fmt.Fprintf(fil, "%s\t%.20f\n", p.Id2PageMap[r], PR[r])
	}
	fmt.Printf("OK!\n")
}
