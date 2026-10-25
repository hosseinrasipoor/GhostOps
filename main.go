package main

import (
	"bufio"
	"flag"
	"fmt"
	"ghostops/parser"
	"os"
)

func main() {
	outFile := flag.String("o", "", "export json file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: ghostops [-o out.json] <logfile>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("error: enter log file path")
		flag.Usage()
		os.Exit(1)
	}

	path := args[0]
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	total := 0
	valid := 0
	malformed := 0

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		total++
		_, err := parser.ParseLine(sc.Text())
		if err != nil {
			malformed++
			continue
		}
		valid++
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("total lines: %d, valid: %d, malformed: %d\n", total, valid, malformed)
	_ = outFile
}
