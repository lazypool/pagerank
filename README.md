# PageRank in WikiPedia

PageRank value of English-language Wikipedia's concepts, based on citation relationship.

What follows is an introduction to the workflow as concisely as possible.

If you want to read the results analysis first, you can skip this part.

## Parse Dump Files

All the code file about this section are stored under `./parse`.

### prepare the dump files

The repository doesn't save WikiPedia's dump files so that manual download is needed to run the program.

All the dump files can be found at [Index of /enwiki/](https://dumps.wikimedia.org/enwiki).

For the requirement of the program, the file storage needs to follow the following conventions:

1. The xml file and the index file must correspond one-to-one and have the same number.
2. All files are in .bz2 compressed format and under the `./enwiki` directory.
3. The index file naming format is like 'index1.txt.bz2', 'index2.txt.bz2'...
4. The xml file naming format is like 'multistream1.xml.bz2', 'multistream2.xml.bz2'...
5. The program has 9 pairs of xml and index by default. If you have more files, please modify the global variable `DUMPNUM` in `./parse/init.go`.

```bash
$ ls ./enwiki
index1.txt.bz2  index4.txt.bz2  index7.txt.bz2  multistream1.xml.bz2 ...
index2.txt.bz2  index5.txt.bz2  index8.txt.bz2  multistream2.xml.bz2 ...
index3.txt.bz2  index6.txt.bz2  index9.txt.bz2  multistream3.xml.bz2 ...
```

### read into the memory

Read the compressed file with the package `compress/bzip2`, `encoding/xml`.

The 3 function `readOffsets`, `readPagesAndIds` and `readPageLinks` are implemented.

- **readOffsets**: read offsets from the given index file, return a slice of offsets, which is used to split the xml file.
- **readPagesAndIds**: read the page2id and id2page map and write it to the given maps, which is used to map page to its id.
- **readPageLinks**: read the page's outlinks from the given file and return a map, wich is used to extract the links.

### find the citations

The function `findLinks` can find all links from the given wiki articale.

It's implemented based on the Regexp.

```Go
commRE = regexp.MustCompile(`(?ms)<!--.*-->`)
nowiRE = regexp.MustCompile(`(?ms)<nowiki>.*</nowiki>`)
linkRE = regexp.MustCompile(`\[\[([^\[\]]+)\]\]`)
```

First replace all the strings match to the commRE and nowRE to the ''.

Then find all strings that match to th linkRE, which are what we want.

### multithreading with goroutine

Implement the multithreading with Go's goroutine and channel mechanism.

The `ParseWiki()` function creates the goroutine's based on the number of dump files and then listens at the channel in the event loop.

The `parseWikiHandler()` is the function that actually conducts the parsing task and put the results to channel if finished.

Using multithreading, the dump files with a total size of about 2G can be processed within only 5 minutes.

## Store Directed Graph

The results of parsing wiki are stored as a directed graph represented by an adjax matrix, which is saved as a json file `adjax.json`.

For matrix storage, the sparse matrix in COO format is used, which is defined in file `./calcu/matrix.go`.

```Go
// the parse COO matrix
type Matrix struct {
	Row []int
	Col []int
	Val []float64
}
```

The file `matrix.go` also realized the many methods about the matrix such as `load a matrix`, `dump a matrix`, `multiply with a vector` and `do a normalization`.

**Note:** 
- The repository doesn't save the `adjax.json` for its too big size. 
- Everytime the program starts, the `ParseWiki()` will run if no `adjax.json` exists.

## Conduct PageRank

The detailed implementation of PageRank algorithm is in `main.go`, which is all encapsulated into the function `PageRank()`.

This function accepts 3 parameters: 

- **pr**: `\*map[int]float64`, the initial PR value, represented as a sparse vector instead of a simple slice.
- **mt**: `\*c.Matrix`, the sparse adjax matrix after normalization by row.
- **n**: `int`, the number of non-zero elements in the PR vector.

```Go
func PageRank(pr *map[int]float64, mt *c.Matrix, n int) {
	var maxC int = 100
	var alpha float64 = 0.85
	var e float64 = (1 - alpha) / float64(N)
	var t float64 = 1e-3

	cp := make(map[int]float64)
	for i := 0; i < maxC; i++ {
		for k, v := range *pr {
			cp[k] = v
		}

		*pr = mt.Mult(pr, true)
		c.VecMult(pr, alpha)
		c.VecAdd(pr, e)

		diff := c.VecComp(pr, &cp)
		if diff < t {
			return
		}
	}
}
```

The function introduces 2 parameters: $\alpha=0.85$ and $e=(1-\alpha)\frac{1}{n}$.

Then it enters a cycle that loops up to 100 cycles.

In this cycle, it update the PR vector with the expressition below again and again.

$$P_{R} \leftarrow \alpha M^T P_{R} + e I, I = (1,1,...,1) \in R^n$$

Before the cycle ends, exit if the difference between the PR vector and its history state is lower than the threshould **t**.

## Results Analysis

All the results of the PageRank are saved in the rank.txt.

It contains 1237934 pages in total, each of which is followed by its PageRank value.

Now let's check the first 20 pages.

```bash
$ head -n20 rank.txt
United States	0.00123069105700100560
Category:United States	0.00108862091398439272
United Kingdom	0.00069379296561897534
The New York Times	0.00068107905080553677
Race and ethnicity in the United States census	0.00067026489903540000
United States Census Bureau	0.00065713042797587015
Category:Living people	0.00058314586064191643
Category:Humans	0.00055703242389706779
Category:Republics	0.00055290078750491669
Category:People by status	0.00050222666948476584
Category:People	0.00045292795468287344
World War II	0.00043457391020196783
England	0.00042651552370922354
Germany	0.00040792081123312704
New York City	0.00040091665642280394
English	0.00039417570757724579
Spain	0.00036951325454533893
London	0.00035584372074505755
Country	0.00034945661752211124
Latin	0.00033917677688987232
```

Based on the list provided, it seems that the most important page is **United States** with a PageRank of about 1.23e-3. The second most important page is **Category:United States** with a PageRank of about 1.08e-3, which is a bit repetitive with the first one. The third most important page is **United Kingdom** with a PageRank of about 0.69e-3.

From this point of view, the concepts of vagueness and generality are even more important. Because they are often referenced by important concepts. This is also the case of **Category:Living people**, **Category:Humans**, **Category:Republics**, **Category:People by status** and **Category: People**, etc.

On the other hand, it's interesting to note that **The New York Times** and **World war II** also obtain a relatively high scores, which proves that things or iconic events that have a greater impact on society and the world can also be important concepts.

In addition to the above two situations, place names are often regarded as important concepts, such as the **England**, **Germany**, **Spain**, **London**, which somehow reflects the importance of each region in the English context.
