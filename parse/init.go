package parse

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// the num of the wiki dump
const DUMPNUM int = 9

// the filename of xml and index
var XMLFILES [DUMPNUM]string
var INDEXFILES [DUMPNUM]string

// the prefix of progress bars
var PREFIXS [DUMPNUM]string

// the channel protoc
type cRES struct {
	UID   int
	PAGES map[int][]int
}

// the parsewiki channel
var workLoad = make(chan [2]int)
var handGain = make(chan *cRES)

// the page & id map
var Page2IdMap map[string]int
var Id2PageMap map[int]string

// the findlinks regexp
var commRE *regexp.Regexp
var linkRE *regexp.Regexp
var nowiRE *regexp.Regexp

func init() {
	// generating wiki dump file paths
	fmt.Printf("Generating Wiki dump file paths...")
	for i := 0; i < DUMPNUM; i++ {
		XMLFILES[i] = "enwiki/multistream" + strconv.Itoa(i+1) + ".xml.bz2"
		INDEXFILES[i] = "enwiki/index" + strconv.Itoa(i+1) + ".txt.bz2"
		PREFIXS[i] = "multistream-" + strconv.Itoa(i+1) + " "
	}
	fmt.Printf("OK!\n\n")

	// Check if the file existing
	for i := 0; i < DUMPNUM; i++ {
		if _, err := os.Stat(XMLFILES[i]); err != nil {
			fmt.Println("Can't open file. ", err)
			return
		}
		if _, err := os.Stat(INDEXFILES[i]); err != nil {
			fmt.Println("Can't open file. ", err)
			return
		}
	}

	// initialize the regexp
	commRE = regexp.MustCompile(`(?ms)<!--.*-->`)
	nowiRE = regexp.MustCompile(`(?ms)<nowiki>.*</nowiki>`)
	linkRE = regexp.MustCompile(`\[\[([^\[\]]+)\]\]`)

	// reading the pages & ids map
	fmt.Printf("Reading Page & Id Map...")
	Page2IdMap = make(map[string]int)
	Id2PageMap = make(map[int]string)
	for i := 0; i < DUMPNUM; i++ {
		readPagesAndIds(INDEXFILES[i], &Page2IdMap, &Id2PageMap)
	}
	fmt.Printf("OK!\n\n")
}
