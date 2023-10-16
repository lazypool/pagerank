package parse

import (
	"strings"
)

// find all links from the given wiki articale
func findLinks(txt string) []int {
	// match all links with regex
	clr := nowiRE.ReplaceAllString(commRE.ReplaceAllString(txt, ""), "")
	fnd := linkRE.FindAllStringSubmatch(clr, -1)

	// for saving the link's id
	var res []int

	// go through found links
	for _, exp := range fnd {
		// the links may between '|'
		for _, l := range strings.Split(exp[1], "|") {
			// insert the link's id if existed
			if id := Page2IdMap[l]; id != 0 {
				res = append(res, id)
			}
		}
	}
	return res
}
