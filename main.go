package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var quietMode bool
	var dontAppend bool
	flag.BoolVar(&quietMode, "q", false, "quiet mode (no output at all)")
	flag.BoolVar(&dontAppend, "da", false, "don't append ")
	
	flag.Parse()

	fn := flag.Arg(0)

	lines := make(map[string]bool)

	var f io.WriteCloser

	if fn != "" {
		// read the whole file into a map if it exists
		r, err := os.Open(fn)
		if err == nil {
			sc := bufio.NewScanner(r)

			for sc.Scan() {
				lines[sc.Text()] = true
			}
			r.Close()
		}

		if !dontAppend {
			// re-open the file for appending new stuff
			f, err = os.OpenFile(fn, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open file for writing: %s\n", err)
				return
			}
			defer f.Close()
		}

	}

	// read the lines, append and output them if they're new
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		line := sc.Text()
		if lines[line] {
			continue
		}

		// add the line to the map so we don't get any duplicates from stdin
		lines[line] = true

		if !quietMode {
			fmt.Println(line)
		}

		if fn != "" && !dontAppend {
			fmt.Fprintf(f, "%s\n", line)
		}
	}
}
