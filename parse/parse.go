package parse

import (
	"fmt"
	"os"
	c "pagerank/calcu"
	u "pagerank/utils"
	"time"
)

// parse all the wiki dump to a matrix then dump the matrix
func ParseWiki() {
	// record the parse time cost
	fmt.Printf("Parsing the Wiki dumps...\n")
	time1 := time.Now()
	defer func() {
		time2 := time.Now()
		cost := time2.Sub(time1).Seconds()
		fmt.Printf("Program ParseWiki exited. Cost %.2fs\n\n", cost)
	}()

	// initialize variables
	var MT c.Matrix
	var UnFin int = 0
	var Time = time.Now()
	var Loop bool = true
	var UI bool = false

	// use to manage the progress bars
	dones := make([]int, DUMPNUM)
	loads := make([]int, DUMPNUM)

	// traverse xml files and create goroutine
	for i := 0; i < DUMPNUM; i++ {
		fmt.Printf("\n")
		uid := int(i)
		go func() {
			parseWikiHandler(XMLFILES[uid], INDEXFILES[uid], uid)
		}()
	}

	// listen at the channel
	for Loop {
		select {
		case c := <-workLoad:
			// store the workload and add a unfinished task
			loads[c[0]] = c[1]
			// start print progress if all task started
			if UnFin++; UnFin > DUMPNUM-1 {
				UI = true
			}
		case r := <-handGain:
			// store the work result
			for src, jps := range r.PAGES {
				for _, dst := range jps {
					MT.Set(src, dst, 1)
				}
			}
			// then num of done works increase
			if dones[r.UID]++; !(dones[r.UID] < loads[r.UID]) {
				if UnFin--; UnFin < 1 {
					// quit if there is no unfinished task
					Loop = false
				}
			}
		default:
			// print progress with a debounce timer
			func() {
				curr := time.Now()
				if curr.Sub(Time).Seconds() > 0.5 && UI {
					Time = time.Now()
					u.ProgressBar(DUMPNUM, dones, loads, PREFIXS[:])
				}
			}()
		}
	}

	// save the matrix to json file
	MT.Dump("adjax.json")
}

// read page2outlinks map from given xml file and index file then push into channel
func parseWikiHandler(xpath string, ipath string, uid int) {
	// read the offsets
	offs := readOffsets(ipath)

	// report the workload
	load := len(offs)
	workLoad <- [2]int{uid, load}

	// the last byte of the file read
	fil, _ := os.Stat(xpath)
	offs = append(offs, int(fil.Size()))

	// read the page2outlinks and report it by channel
	for i := 0; i < load; i++ {
		off := offs[i]
		lim := offs[i+1] - off
		pages := readPageLinks(xpath, off, lim)
		handGain <- &cRES{UID: uid, PAGES: pages}
	}
}
