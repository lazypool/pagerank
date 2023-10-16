package parse

import (
	"bufio"
	"compress/bzip2"
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"strings"
)

// read offsets from the given file, return a slice of offsets
func readOffsets(path string) []int {
	fil, _ := os.Open(path)
	defer fil.Close()

	raw := bzip2.NewReader(fil)
	rdr := bufio.NewReader(raw)

	var offs []int
	var curr int
	for {
		// in each line
		line, err := rdr.ReadString('\n')
		if err == io.EOF {
			break
		}

		// the offset is the first element
		off := strings.SplitN(line, ":", 3)[0]
		num, _ := strconv.Atoi(off)

		// only insert it if it's a new unseen offset
		if num != curr {
			offs = append(offs, num)
			curr = num
		}
	}

	return offs
}

// read the page2id and id2page map and write it to the given maps
func readPagesAndIds(path string, p2i *map[string]int, i2p *map[int]string) {
	fil, _ := os.Open(path)
	defer fil.Close()

	raw := bzip2.NewReader(fil)
	rdr := bufio.NewReader(raw)

	for {
		// in each line
		line, err := rdr.ReadString('\n')
		if err == io.EOF {
			break
		}

		// id is the second one and page is the third one
		cuts := strings.SplitN(line, ":", 3)
		id, _ := strconv.Atoi(cuts[1])
		page := strings.TrimRight(cuts[2], "\n")

		// write it to the given maps
		(*p2i)[page] = id
		(*i2p)[id] = page
	}
}

// read the page2outlinks map from the given file and return it
func readPageLinks(path string, off int, lim int) map[int][]int {
	fil, _ := os.Open(path)
	defer fil.Close()

	// obtain the decoder by given 'path', 'offset' and 'limit'
	sec := io.NewSectionReader(fil, int64(off), int64(lim))
	raw := bzip2.NewReader(sec)
	dec := xml.NewDecoder(raw)

	// label, caption, text
	lbl := ""
	cpt := ""
	txt := ""

	res := make(map[int][]int)
	// traverse the xml file
	for {
		token, err := dec.Token()
		if err != nil {
			break
		}

		// get the current label (e.g. <page>, <title>, <text>)
		if cur, ok := token.(xml.StartElement); ok {
			lbl = cur.Name.Local
		} else if _, ok := token.(xml.EndElement); ok {
			lbl = "EOF"
		}

		// get the current title and current text
		if cur, ok := token.(xml.CharData); ok {
			switch lbl {
			case "title":
				cpt = string(cur)
			case "text":
				txt = string(cur)
			}
		}

		// if all title and text found
		if cpt != "" && txt != "" {
			// find the source pageId and the destination pageId
			src := Page2IdMap[cpt]
			jps := findLinks(txt)
			// append it only if found and outlinks is not empty
			if src != 0 && len(jps) > 0 {
				res[src] = jps
			}
			// clear to find new title and text
			cpt, txt = "", ""
		}
	}

	return res
}
