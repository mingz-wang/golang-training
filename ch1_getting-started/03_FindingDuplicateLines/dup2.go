// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	line2file := make(map[string][]string)

	files := os.Args[1:]

	if len(files) == 0 {
		// countLines(os.Stdin, counts)
	} else {
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, line2file)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\t%s\n", n, line, line2file[line])
		}
	}
}

func countLines(f *os.File, counts map[string]int, line2file map[string][]string) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
		line2file[input.Text()] = append(line2file[input.Text()], f.Name())
	}
	// NOTE: ignoring potential errors from input.Err()
}
