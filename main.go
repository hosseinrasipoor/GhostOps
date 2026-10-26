package main

import (
	"bufio"
	"flag"
	"fmt"
	"ghostops/analyzer"
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

	an := analyzer.NewAnalyzer()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		entry, err := parser.ParseLine(sc.Text())
		if err != nil {
			an.AddBad()
			continue
		}
		an.AddGood(entry)
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("error reading file: %v\n", err)
		os.Exit(1)
	}

	res := an.GetResult()
	fmt.Printf("total lines: %d, valid: %d, malformed: %d\n", res.TotalLines, res.ValidEntries, res.MalformedLines)
	fmt.Printf("incidents: %d\n", len(res.Incidents))
	_ = outFile
}
